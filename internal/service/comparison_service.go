package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/domain/repository"
	"paper_ai/internal/infrastructure/comparison"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
)

// ComparisonService 对比服务
type ComparisonService struct {
	polishRepo       repository.PolishRepository
	diffEngine       *comparison.DiffEngine
	positionCalc     *comparison.PositionCalculator
	classifier       *comparison.ChangeClassifier
	reasonGenerator  *comparison.ReasonGenerator
}

// NewComparisonService 创建对比服务
func NewComparisonService(polishRepo repository.PolishRepository) *ComparisonService {
	return &ComparisonService{
		polishRepo:      polishRepo,
		diffEngine:      comparison.NewDiffEngine(),
		positionCalc:    comparison.NewPositionCalculator(),
		classifier:      comparison.NewChangeClassifier(),
		reasonGenerator: comparison.NewReasonGenerator(),
	}
}

// GenerateComparison 生成对比数据
func (s *ComparisonService) GenerateComparison(ctx context.Context, traceID string) (*model.ComparisonResult, error) {
	// 1. 获取润色记录
	record, err := s.polishRepo.GetByTraceID(ctx, traceID)
	if err != nil {
		logger.Error("failed to get polish record", zap.String("trace_id", traceID), zap.Error(err))
		return nil, apperrors.NewNotFoundError("润色记录不存在")
	}

	// 2. 如果已有对比数据，直接返回
	if record.ComparisonData != "" {
		return s.parseComparisonData(record)
	}

	// 3. 生成对比数据
	result, err := s.generateComparisonData(record)
	if err != nil {
		logger.Error("failed to generate comparison data", zap.String("trace_id", traceID), zap.Error(err))
		return nil, err
	}

	// 4. 保存对比数据
	if err := s.saveComparisonData(ctx, record, result); err != nil {
		logger.Warn("failed to save comparison data", zap.String("trace_id", traceID), zap.Error(err))
		// 不返回错误，继续返回生成的对比数据
	}

	return result, nil
}

// generateComparisonData 生成对比数据
func (s *ComparisonService) generateComparisonData(record *entity.PolishRecord) (*model.ComparisonResult, error) {
	original := record.OriginalContent
	polished := record.PolishedContent

	// 1. 运行 diff 算法
	diffs := s.diffEngine.GenerateDiff(original, polished)

	// 2. 提取修改信息
	changes := s.diffEngine.GetChanges(diffs)

	// 3. 计算位置
	positions := s.positionCalc.CalculatePositions(polished, changes)

	// 4. 生成标注列表
	annotations := s.buildAnnotations(positions)

	// 5. 计算元数据和统计信息
	metadata, statistics := s.calculateStats(original, polished, annotations)

	return &model.ComparisonResult{
		TraceID:         record.TraceID,
		OriginalContent: original,
		PolishedContent: polished,
		Annotations:     annotations,
		Metadata:        metadata,
		Statistics:      statistics,
	}, nil
}

// buildAnnotations 构建标注列表
func (s *ComparisonService) buildAnnotations(positions []comparison.PositionInfo) []model.Change {
	annotations := make([]model.Change, 0, len(positions))

	for i, pos := range positions {
		// 分类修改类型
		changeType := s.classifier.Classify(pos.OriginalText, pos.PolishedText)

		// 生成修改理由
		reason := s.reasonGenerator.Generate(changeType, pos.OriginalText, pos.PolishedText)

		// 生成替代方案
		alternatives := s.reasonGenerator.GenerateAlternatives(changeType, pos.OriginalText)

		// 计算置信度
		confidence := s.reasonGenerator.CalculateConfidence(changeType, pos.OriginalText, pos.PolishedText)

		// 获取影响维度
		impact := s.reasonGenerator.GetImpact(changeType)

		// 建议高亮颜色
		highlightColor := s.classifier.SuggestHighlightColor(changeType)

		annotations = append(annotations, model.Change{
			ID:   fmt.Sprintf("change_%d", i+1),
			Type: changeType,
			PolishedPosition: model.Position{
				Start: pos.Start,
				End:   pos.End,
				Line:  pos.Line,
			},
			PolishedText:   pos.PolishedText,
			OriginalText:   pos.OriginalText,
			Reason:         reason,
			Alternatives:   alternatives,
			Confidence:     confidence,
			Impact:         impact,
			HighlightColor: highlightColor,
			Status:         model.ActionStatusPending,
		})
	}

	return annotations
}

// calculateStats 计算统计信息
func (s *ComparisonService) calculateStats(original, polished string, annotations []model.Change) (model.Metadata, model.Statistics) {
	originalWordCount := comparison.CountWords(original)
	polishedWordCount := comparison.CountWords(polished)

	// 统计各类修改数量
	var vocabCount, grammarCount, structureCount int
	for _, ann := range annotations {
		switch ann.Type {
		case model.ChangeTypeVocabulary:
			vocabCount++
		case model.ChangeTypeGrammar:
			grammarCount++
		case model.ChangeTypeStructure:
			structureCount++
		}
	}

	// 计算学术性提升（简单算法：词汇优化占比 * 100）
	academicImprovement := 0.0
	if len(annotations) > 0 {
		academicImprovement = float64(vocabCount) / float64(len(annotations)) * 100
		// 限制在0-100之间
		if academicImprovement > 100 {
			academicImprovement = 100
		}
	}

	metadata := model.Metadata{
		OriginalWordCount:        originalWordCount,
		PolishedWordCount:        polishedWordCount,
		TotalChanges:             len(annotations),
		AcademicScoreImprovement: academicImprovement,
	}

	statistics := model.Statistics{
		VocabularyChanges: vocabCount,
		GrammarChanges:    grammarCount,
		StructureChanges:  structureCount,
	}

	return metadata, statistics
}

// saveComparisonData 保存对比数据
func (s *ComparisonService) saveComparisonData(ctx context.Context, record *entity.PolishRecord, result *model.ComparisonResult) error {
	// 序列化对比数据
	dataJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal comparison data: %w", err)
	}

	// 收集已接受和已拒绝的修改ID列表
	acceptedIDs := []string{}
	rejectedIDs := []string{}
	for _, ann := range result.Annotations {
		if ann.Status == model.ActionStatusAccepted {
			acceptedIDs = append(acceptedIDs, ann.ID)
		} else if ann.Status == model.ActionStatusRejected {
			rejectedIDs = append(rejectedIDs, ann.ID)
		}
	}

	// 更新记录
	record.ComparisonData = string(dataJSON)
	record.ChangesCount = result.Metadata.TotalChanges
	record.AcceptedChanges = acceptedIDs
	record.RejectedChanges = rejectedIDs
	record.UpdatedAt = time.Now()

	return s.polishRepo.Update(ctx, record)
}

// parseComparisonData 解析对比数据
func (s *ComparisonService) parseComparisonData(record *entity.PolishRecord) (*model.ComparisonResult, error) {
	var result model.ComparisonResult
	if err := json.Unmarshal([]byte(record.ComparisonData), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal comparison data: %w", err)
	}
	return &result, nil
}

// GetComparison 获取对比数据（外部接口）
func (s *ComparisonService) GetComparison(ctx context.Context, traceID string, userID int64) (*model.ComparisonResult, error) {
	// 1. 获取润色记录
	record, err := s.polishRepo.GetByTraceID(ctx, traceID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("润色记录不存在")
	}

	// 2. 验证权限
	if record.UserID != userID {
		return nil, apperrors.NewForbiddenError("无权访问该记录")
	}

	// 3. 生成或获取对比数据
	return s.GenerateComparison(ctx, traceID)
}

// ApplyAction 应用用户操作（接受/拒绝修改）
func (s *ComparisonService) ApplyAction(ctx context.Context, traceID string, userID int64, req *model.ChangeActionRequest) (*model.ChangeActionResponse, error) {
	// 1. 获取记录
	record, err := s.polishRepo.GetByTraceID(ctx, traceID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("润色记录不存在")
	}

	// 2. 验证权限
	if record.UserID != userID {
		return nil, apperrors.NewForbiddenError("无权访问该记录")
	}

	// 3. 获取对比数据
	result, err := s.GenerateComparison(ctx, traceID)
	if err != nil {
		return nil, err
	}

	// 4. 查找目标修改
	changeIndex := -1
	for i, ann := range result.Annotations {
		if ann.ID == req.ChangeID {
			changeIndex = i
			break
		}
	}

	if changeIndex == -1 {
		return nil, apperrors.NewInvalidParameterError("修改不存在")
	}

	// 5. 应用操作
	if req.Action == "accept" {
		result.Annotations[changeIndex].Status = model.ActionStatusAccepted
	} else if req.Action == "reject" {
		result.Annotations[changeIndex].Status = model.ActionStatusRejected
	}

	// 6. 生成更新后的内容
	updatedContent := s.applyChanges(result)

	// 7. 统计已应用和待处理的修改
	appliedChanges := []string{req.ChangeID}
	pendingChanges := []string{}
	for _, ann := range result.Annotations {
		if ann.Status == model.ActionStatusPending {
			pendingChanges = append(pendingChanges, ann.ID)
		}
	}

	// 8. 保存更新后的对比数据和最终文本
	record.FinalContent = updatedContent // 保存最终文本
	if err := s.saveComparisonData(ctx, record, result); err != nil {
		logger.Warn("failed to save updated comparison data", zap.Error(err))
	}

	return &model.ChangeActionResponse{
		Success:        true,
		UpdatedContent: updatedContent,
		AppliedChanges: appliedChanges,
		PendingChanges: pendingChanges,
	}, nil
}

// applyChanges 应用所有接受的修改，生成最终文本
// 策略：从原文开始，应用所有 accepted 状态的修改
func (s *ComparisonService) applyChanges(result *model.ComparisonResult) string {
	// 1. 如果没有任何修改或全部拒绝，返回原文
	hasAcceptedChanges := false
	for _, ann := range result.Annotations {
		if ann.Status == model.ActionStatusAccepted {
			hasAcceptedChanges = true
			break
		}
	}

	if !hasAcceptedChanges {
		return result.OriginalContent
	}

	// 2. 如果全部接受，返回润色后文本
	allAccepted := true
	for _, ann := range result.Annotations {
		if ann.Status != model.ActionStatusAccepted {
			allAccepted = false
			break
		}
	}

	if allAccepted {
		return result.PolishedContent
	}

	// 3. 部分接受：需要重新构建文本
	// 使用 DiffEngine 重新计算，基于原文应用接受的修改
	return s.rebuildTextWithAcceptedChanges(result)
}

// rebuildTextWithAcceptedChanges 重新构建文本（只应用接受的修改）
func (s *ComparisonService) rebuildTextWithAcceptedChanges(result *model.ComparisonResult) string {
	// 简化实现：使用字符串替换
	// 从原文开始，逐个应用接受的修改

	finalText := result.OriginalContent

	// 收集所有接受的修改，并按原文中的位置排序（从后往前替换，避免位置偏移）
	type acceptedChange struct {
		originalText string
		polishedText string
		firstIndex   int // 在原文中首次出现的位置
	}

	acceptedChanges := []acceptedChange{}

	for _, ann := range result.Annotations {
		if ann.Status == model.ActionStatusAccepted {
			// 在原文中查找该词的位置
			idx := strings.Index(finalText, ann.OriginalText)
			if idx != -1 {
				acceptedChanges = append(acceptedChanges, acceptedChange{
					originalText: ann.OriginalText,
					polishedText: ann.PolishedText,
					firstIndex:   idx,
				})
			}
		}
	}

	// 按位置从后往前排序（避免替换后位置偏移）
	sort.Slice(acceptedChanges, func(i, j int) bool {
		return acceptedChanges[i].firstIndex > acceptedChanges[j].firstIndex
	})

	// 应用替换
	for _, change := range acceptedChanges {
		// 只替换第一次出现的位置
		idx := strings.Index(finalText, change.originalText)
		if idx != -1 {
			finalText = finalText[:idx] + change.polishedText + finalText[idx+len(change.originalText):]
		}
	}

	return finalText
}

// BatchAcceptAll 一键接受所有修改
func (s *ComparisonService) BatchAcceptAll(ctx context.Context, traceID string, userID int64) (*model.BatchActionResponse, error) {
	// 1. 获取记录
	record, err := s.polishRepo.GetByTraceID(ctx, traceID)
	if err != nil {
		logger.Error("failed to get polish record", zap.String("trace_id", traceID), zap.Error(err))
		return nil, apperrors.NewNotFoundError("润色记录不存在")
	}

	// 2. 验证权限
	if record.UserID != userID {
		return nil, apperrors.NewForbiddenError("无权访问该记录")
	}

	// 3. 获取对比数据
	result, err := s.GenerateComparison(ctx, traceID)
	if err != nil {
		return nil, err
	}

	// 4. 批量接受所有修改
	appliedCount := 0
	for i := range result.Annotations {
		if result.Annotations[i].Status == model.ActionStatusPending {
			result.Annotations[i].Status = model.ActionStatusAccepted
			appliedCount++
		}
	}

	// 5. 生成最终文本（全部接受 = 润色后文本）
	record.FinalContent = result.PolishedContent

	// 6. 保存更新后的对比数据
	if err := s.saveComparisonData(ctx, record, result); err != nil {
		logger.Warn("failed to save updated comparison data", zap.Error(err))
	}

	return &model.BatchActionResponse{
		Success:        true,
		UpdatedContent: result.PolishedContent,
		AppliedCount:   appliedCount,
	}, nil
}
