package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"paper_ai/internal/config"
	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/domain/repository"
	"paper_ai/internal/infrastructure/ai"
	"paper_ai/internal/infrastructure/ai/types"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/idgen"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
)

// PolishMultiVersionService 多版本润色服务
type PolishMultiVersionService struct {
	providerFactory *ai.ProviderFactory
	polishRepo      repository.PolishRepository
	versionRepo     repository.PolishVersionRepository
	promptService   *PromptService
	featureService  *FeatureService
}

// NewPolishMultiVersionService 创建多版本润色服务
func NewPolishMultiVersionService(
	factory *ai.ProviderFactory,
	polishRepo repository.PolishRepository,
	versionRepo repository.PolishVersionRepository,
	promptService *PromptService,
	featureService *FeatureService,
) *PolishMultiVersionService {
	return &PolishMultiVersionService{
		providerFactory: factory,
		polishRepo:      polishRepo,
		versionRepo:     versionRepo,
		promptService:   promptService,
		featureService:  featureService,
	}
}

// PolishMultiVersion 执行多版本润色
func (s *PolishMultiVersionService) PolishMultiVersion(ctx context.Context, req *model.PolishMultiVersionRequest, userID int64) (*model.PolishMultiVersionResponse, error) {
	startTime := time.Now()

	// 生成TraceID
	traceID, err := s.generateTraceID(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("multi-version polish started",
		zap.String("trace_id", traceID),
		zap.Int64("user_id", userID),
		zap.String("language", req.Language),
		zap.String("style", req.Style))

	// 1. 权限检查
	hasPermission, reason, err := s.featureService.CheckMultiVersionPermission(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		logger.Warn("user does not have multi-version permission", zap.Int64("user_id", userID), zap.String("reason", reason))
		return nil, apperrors.NewForbiddenError(reason)
	}

	// 2. 参数验证和设置默认值
	if err := s.validateAndSetDefaults(req); err != nil {
		return nil, apperrors.NewInvalidParameterError(err.Error())
	}

	// 3. 获取AI提供商
	provider, err := s.getProvider(req.Provider)
	if err != nil {
		return nil, err
	}

	// 4. 创建主记录
	mainRecord := &entity.PolishRecord{
		TraceID:         traceID,
		UserID:          userID,
		OriginalContent: req.Content,
		Style:           req.Style,
		Language:        req.Language,
		OriginalLength:  len(req.Content),
		Provider:        req.Provider,
		Mode:            entity.ModeMulti,
		Status:          "processing",
	}

	if err := s.polishRepo.Create(ctx, mainRecord); err != nil {
		logger.Error("failed to create main record", zap.Error(err))
		return nil, fmt.Errorf("failed to create main record: %w", err)
	}

	// 5. 确定要生成的版本类型
	versionTypes := s.determineVersionTypes(req.Versions)

	// 6. 并发调用AI生成多个版本
	versionResults := s.generateVersionsConcurrently(ctx, versionTypes, req, provider, mainRecord.ID)

	// 7. 统计结果
	successCount := 0
	failedCount := 0
	totalProcessTime := 0
	versions := make(map[string]*model.VersionResult)

	for versionType, result := range versionResults {
		versions[versionType] = result
		totalProcessTime += result.ProcessTimeMs

		if result.Status == "success" {
			successCount++
		} else {
			failedCount++
		}
	}

	// 8. 更新主记录状态
	status := "success"
	if failedCount == len(versionTypes) {
		status = "failed"
	} else if failedCount > 0 {
		status = "partial"
	}

	// 使用第一个成功版本的内容作为主记录的 PolishedContent（如果有的话）
	polishedContent := ""
	polishedLength := 0
	modelUsed := ""
	for _, result := range versions {
		if result.Status == "success" {
			polishedContent = result.PolishedContent
			polishedLength = result.PolishedLength
			modelUsed = result.ModelUsed
			break
		}
	}

	mainRecord.Status = status
	mainRecord.PolishedContent = polishedContent
	mainRecord.PolishedLength = polishedLength
	mainRecord.Model = modelUsed
	mainRecord.ProcessTimeMs = totalProcessTime

	if err := s.polishRepo.Update(ctx, mainRecord); err != nil {
		logger.Error("failed to update main record", zap.Error(err))
		// 不返回错误，因为主要工作已完成
	}

	// 记录总耗时
	totalElapsed := time.Since(startTime).Milliseconds()
	logger.Info("multi-version polish completed",
		zap.String("trace_id", traceID),
		zap.Int("success_count", successCount),
		zap.Int("failed_count", failedCount),
		zap.Int64("total_elapsed_ms", totalElapsed))

	// 9. 构建响应
	return &model.PolishMultiVersionResponse{
		TraceID:        traceID,
		OriginalLength: len(req.Content),
		Versions:       versions,
		ProviderUsed:   req.Provider,
	}, nil
}

// generateVersionsConcurrently 并发生成多个版本
func (s *PolishMultiVersionService) generateVersionsConcurrently(
	ctx context.Context,
	versionTypes []string,
	req *model.PolishMultiVersionRequest,
	provider ai.AIProvider,
	recordID int64,
) map[string]*model.VersionResult {
	var wg sync.WaitGroup
	results := make(map[string]*model.VersionResult)
	mu := sync.Mutex{}

	// 并发生成每个版本
	for _, versionType := range versionTypes {
		wg.Add(1)
		go func(vt string) {
			defer wg.Done()

			result := s.generateSingleVersion(ctx, vt, req, provider, recordID)

			mu.Lock()
			results[vt] = result
			mu.Unlock()
		}(versionType)
	}

	// 等待所有版本完成
	wg.Wait()

	return results
}

// generateSingleVersion 生成单个版本
func (s *PolishMultiVersionService) generateSingleVersion(
	ctx context.Context,
	versionType string,
	req *model.PolishMultiVersionRequest,
	provider ai.AIProvider,
	recordID int64,
) *model.VersionResult {
	startTime := time.Now()

	logger.Info("generating version",
		zap.String("version_type", versionType),
		zap.String("language", req.Language),
		zap.String("style", req.Style))

	// 1. 渲染Prompt
	renderedPrompt, err := s.promptService.RenderPrompt(ctx, versionType, req.Language, req.Style, req.Content)
	if err != nil {
		logger.Error("failed to render prompt",
			zap.String("version_type", versionType),
			zap.Error(err))
		return &model.VersionResult{
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("failed to render prompt: %v", err),
		}
	}

	// 2. 调用AI
	polishReq := &types.PolishRequest{
		Content:  renderedPrompt.UserPrompt,
		Style:    req.Style,
		Language: req.Language,
	}

	polishResp, err := provider.Polish(ctx, polishReq)
	if err != nil {
		logger.Error("failed to call AI provider",
			zap.String("version_type", versionType),
			zap.Error(err))

		// 保存失败记录
		s.saveFailedVersion(ctx, recordID, versionType, renderedPrompt.PromptID, err)

		return &model.VersionResult{
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("AI call failed: %v", err),
		}
	}

	processTimeMs := int(time.Since(startTime).Milliseconds())

	// 3. 保存版本记录
	version := &entity.PolishVersion{
		RecordID:        recordID,
		VersionType:     versionType,
		PolishedContent: polishResp.PolishedContent,
		PolishedLength:  len(polishResp.PolishedContent),
		Suggestions:     polishResp.Suggestions,
		ModelUsed:       polishResp.ModelUsed,
		PromptID:        renderedPrompt.PromptID,
		ProcessTimeMs:   processTimeMs,
		Status:          "success",
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		logger.Error("failed to save version record",
			zap.String("version_type", versionType),
			zap.Error(err))
		// 不返回错误，因为AI调用已成功
	}

	// 4. 增加Prompt使用次数
	if err := s.promptService.IncrementUsage(ctx, renderedPrompt.PromptID); err != nil {
		logger.Error("failed to increment prompt usage", zap.Error(err))
		// 不影响主流程
	}

	logger.Info("version generated successfully",
		zap.String("version_type", versionType),
		zap.Int("process_time_ms", processTimeMs))

	return &model.VersionResult{
		PolishedContent: polishResp.PolishedContent,
		PolishedLength:  len(polishResp.PolishedContent),
		Suggestions:     polishResp.Suggestions,
		ProcessTimeMs:   processTimeMs,
		ModelUsed:       polishResp.ModelUsed,
		Status:          "success",
	}
}

// saveFailedVersion 保存失败的版本记录
func (s *PolishMultiVersionService) saveFailedVersion(ctx context.Context, recordID int64, versionType string, promptID int64, err error) {
	version := &entity.PolishVersion{
		RecordID:      recordID,
		VersionType:   versionType,
		PromptID:      promptID,
		Status:        "failed",
		ErrorMessage:  err.Error(),
		ProcessTimeMs: 0,
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		logger.Error("failed to save failed version record", zap.Error(err))
	}
}

// validateAndSetDefaults 验证参数并设置默认值
func (s *PolishMultiVersionService) validateAndSetDefaults(req *model.PolishMultiVersionRequest) error {
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}

	if len(req.Content) > 10000 {
		return fmt.Errorf("content too long (max 10000 characters)")
	}

	// 设置默认值
	if req.Style == "" {
		req.Style = "academic"
	}
	if req.Language == "" {
		req.Language = "en"
	}
	if req.Provider == "" {
		req.Provider = config.Get().AI.DefaultProvider
	}

	return nil
}

// getProvider 获取AI提供商
func (s *PolishMultiVersionService) getProvider(providerName string) (ai.AIProvider, error) {
	if providerName == "" {
		return s.providerFactory.GetDefaultProvider()
	}
	return s.providerFactory.GetProvider(providerName)
}

// determineVersionTypes 确定要生成的版本类型
func (s *PolishMultiVersionService) determineVersionTypes(requestedVersions []string) []string {
	// 如果请求中指定了版本，则使用指定的版本
	if len(requestedVersions) > 0 {
		validVersions := make([]string, 0, len(requestedVersions))
		for _, v := range requestedVersions {
			if entity.IsValidVersionType(v) {
				validVersions = append(validVersions, v)
			}
		}
		if len(validVersions) > 0 {
			return validVersions
		}
	}

	// 否则生成全部3个版本
	return []string{
		entity.VersionTypeConservative,
		entity.VersionTypeBalanced,
		entity.VersionTypeAggressive,
	}
}

// generateTraceID 生成TraceID
func (s *PolishMultiVersionService) generateTraceID(ctx context.Context) (string, error) {
	// 从context中获取traceID，如果没有则生成
	traceID, ok := ctx.Value("trace_id").(string)
	if ok && traceID != "" {
		return traceID, nil
	}

	// 使用Snowflake ID生成器
	id, err := idgen.GenerateID()
	if err != nil {
		logger.Error("failed to generate trace ID", zap.Error(err))
		// 降级方案：使用时间戳
		return strconv.FormatInt(time.Now().UnixNano(), 10), nil
	}

	return strconv.FormatInt(id, 10), nil
}
