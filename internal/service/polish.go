package service

import (
	"context"
	"strconv"
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

// PolishService 润色服务
type PolishService struct {
	providerFactory *ai.ProviderFactory
	polishRepo      repository.PolishRepository // 仓储接口（依赖注入）
}

// NewPolishService 创建润色服务
func NewPolishService(factory *ai.ProviderFactory, repo repository.PolishRepository) *PolishService {
	return &PolishService{
		providerFactory: factory,
		polishRepo:      repo,
	}
}

// Polish 执行段落润色
func (s *PolishService) Polish(ctx context.Context, req *model.PolishRequest, userID int64) (*types.PolishResponse, error) {
	startTime := time.Now()

	// 从context中获取traceID，如果没有则生成唯一ID
	traceID, ok := ctx.Value("trace_id").(string)
	if !ok || traceID == "" {
		// 使用 Snowflake ID 生成器生成纯数字 TraceID
		id, err := idgen.GenerateID()
		if err != nil {
			logger.Error("failed to generate trace ID", zap.Error(err))
			// 降级方案：使用时间戳
			traceID = strconv.FormatInt(time.Now().UnixNano(), 10)
		} else {
			traceID = strconv.FormatInt(id, 10)
		}
	}

	// 参数验证
	if err := req.Validate(); err != nil {
		logger.Warn("invalid polish request", zap.Error(err))
		// 记录失败的请求
		s.saveFailedRecord(ctx, traceID, req, userID, err)
		return nil, apperrors.NewInvalidParameterError(err.Error())
	}

	// 设置默认值
	req.SetDefaults()

	// 获取AI提供商
	var provider ai.AIProvider
	var err error

	if req.Provider == "" {
		// 使用默认提供商
		provider, err = s.providerFactory.GetDefaultProvider()
		if err != nil {
			logger.Error("failed to get default provider", zap.Error(err))
			s.saveFailedRecord(ctx, traceID, req, userID, err)
			return nil, err
		}
		req.Provider = config.Get().AI.DefaultProvider
	} else {
		// 使用指定的提供商
		provider, err = s.providerFactory.GetProvider(req.Provider)
		if err != nil {
			logger.Error("failed to get provider", zap.String("provider", req.Provider), zap.Error(err))
			s.saveFailedRecord(ctx, traceID, req, userID, err)
			return nil, err
		}
	}

	// 构建AI请求
	aiReq := &types.PolishRequest{
		Content:  req.Content,
		Style:    req.Style,
		Language: req.Language,
	}

	// 调用AI服务
	logger.Info("calling ai provider for polish",
		zap.String("provider", req.Provider),
		zap.Int("content_length", len(req.Content)),
		zap.Int64("user_id", userID),
	)

	resp, err := provider.Polish(ctx, aiReq)
	if err != nil {
		logger.Error("ai provider polish failed",
			zap.String("provider", req.Provider),
			zap.Error(err),
		)
		s.saveFailedRecord(ctx, traceID, req, userID, err)
		return nil, err
	}

	// 计算处理时间
	processTime := time.Since(startTime).Milliseconds()

	// 保存成功记录
	s.saveSuccessRecord(ctx, traceID, req, resp, userID, int(processTime))

	// 设置 TraceID 到响应中
	resp.TraceID = traceID

	logger.Info("polish completed successfully",
		zap.String("provider", req.Provider),
		zap.String("trace_id", traceID),
		zap.Int("original_length", resp.OriginalLength),
		zap.Int("polished_length", resp.PolishedLength),
		zap.Int64("process_time_ms", processTime),
		zap.Int64("user_id", userID),
	)

	return resp, nil
}

// saveSuccessRecord 保存成功记录
func (s *PolishService) saveSuccessRecord(ctx context.Context, traceID string, req *model.PolishRequest, resp *types.PolishResponse, userID int64, processTime int) {
	if s.polishRepo == nil {
		return // 如果没有配置数据库，跳过保存
	}

	record := &entity.PolishRecord{
		TraceID:         traceID,
		UserID:          userID,
		OriginalContent: req.Content,
		Style:           req.Style,
		Language:        req.Language,
		PolishedContent: resp.PolishedContent,
		OriginalLength:  resp.OriginalLength,
		PolishedLength:  resp.PolishedLength,
		Provider:        resp.ProviderUsed,
		Model:           resp.ModelUsed,
		ProcessTimeMs:   processTime,
		Status:          "success",
	}

	if err := s.polishRepo.Create(ctx, record); err != nil {
		logger.Error("failed to save polish record", zap.String("trace_id", traceID), zap.Error(err))
	}
}

// saveFailedRecord 保存失败记录
func (s *PolishService) saveFailedRecord(ctx context.Context, traceID string, req *model.PolishRequest, userID int64, err error) {
	if s.polishRepo == nil {
		return
	}

	record := &entity.PolishRecord{
		TraceID:         traceID,
		UserID:          userID,
		OriginalContent: req.Content,
		Style:           req.Style,
		Language:        req.Language,
		Status:          "failed",
		ErrorMessage:    err.Error(),
	}

	if err := s.polishRepo.Create(ctx, record); err != nil {
		logger.Error("failed to save failed polish record", zap.String("trace_id", traceID), zap.Error(err))
	}
}

// GetRecordByTraceID 根据TraceID获取记录
func (s *PolishService) GetRecordByTraceID(ctx context.Context, traceID string, userID int64) (*entity.PolishRecord, error) {
	if s.polishRepo == nil {
		return nil, apperrors.NewInternalError("database not configured", nil)
	}

	record, err := s.polishRepo.GetByTraceID(ctx, traceID)
	if err != nil {
		return nil, err
	}

	// 验证记录所有权（只能查看自己的记录）
	if record.UserID != userID {
		return nil, apperrors.NewForbiddenError("you don't have permission to access this record")
	}

	return record, nil
}

// ListRecords 获取记录列表
func (s *PolishService) ListRecords(ctx context.Context, opts repository.QueryOptions) ([]*entity.PolishRecord, int64, error) {
	if s.polishRepo == nil {
		return nil, 0, apperrors.NewInternalError("database not configured", nil)
	}

	records, err := s.polishRepo.List(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.polishRepo.Count(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetStatistics 获取统计信息
func (s *PolishService) GetStatistics(ctx context.Context, opts repository.StatisticsOptions) (*repository.Statistics, error) {
	if s.polishRepo == nil {
		return nil, apperrors.NewInternalError("database not configured", nil)
	}
	return s.polishRepo.GetStatistics(ctx, opts)
}
