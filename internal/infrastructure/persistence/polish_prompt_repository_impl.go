package persistence

import (
	"context"
	"fmt"

	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/repository"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// polishPromptRepositoryImpl Prompt仓储实现
type polishPromptRepositoryImpl struct {
	db *gorm.DB
}

// NewPolishPromptRepository 创建Prompt仓储实现
func NewPolishPromptRepository(db *gorm.DB) repository.PolishPromptRepository {
	return &polishPromptRepositoryImpl{db: db}
}

// Create 创建Prompt
func (r *polishPromptRepositoryImpl) Create(ctx context.Context, prompt *entity.PolishPrompt) error {
	po := &PolishPromptPO{}
	po.FromEntity(prompt)

	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		logger.Error("failed to create polish prompt", zap.Error(err))
		return fmt.Errorf("failed to create polish prompt: %w", err)
	}

	// 回写ID和时间戳
	prompt.ID = po.ID
	prompt.CreatedAt = po.CreatedAt
	prompt.UpdatedAt = po.UpdatedAt

	return nil
}

// GetByID 根据ID获取Prompt
func (r *polishPromptRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.PolishPrompt, error) {
	var po PolishPromptPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("polish prompt not found: id=%d", id)
		}
		logger.Error("failed to get polish prompt by id", zap.Int64("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get polish prompt: %w", err)
	}

	return po.ToEntity(), nil
}

// GetActive 获取激活的Prompt
// 查询策略（按优先级）：
// 1. 精确匹配：versionType + language + style
// 2. 降级匹配：versionType + language + style='all'
// 3. 再降级：versionType + language='all' + style='all'
func (r *polishPromptRepositoryImpl) GetActive(ctx context.Context, versionType, language, style string) (*entity.PolishPrompt, error) {
	var po PolishPromptPO

	// 策略1：精确匹配
	err := r.db.WithContext(ctx).
		Where("version_type = ? AND language = ? AND style = ? AND is_active = ?", versionType, language, style, true).
		Order("version DESC, weight DESC").
		First(&po).Error

	if err == nil {
		logger.Info("found exact match prompt",
			zap.String("version_type", versionType),
			zap.String("language", language),
			zap.String("style", style),
			zap.Int64("prompt_id", po.ID))
		return po.ToEntity(), nil
	}

	if err != gorm.ErrRecordNotFound {
		logger.Error("failed to query exact match prompt", zap.Error(err))
		return nil, fmt.Errorf("failed to query prompt: %w", err)
	}

	// 策略2：降级匹配（style='all'）
	err = r.db.WithContext(ctx).
		Where("version_type = ? AND language = ? AND style = ? AND is_active = ?", versionType, language, "all", true).
		Order("version DESC, weight DESC").
		First(&po).Error

	if err == nil {
		logger.Info("found style-generic prompt",
			zap.String("version_type", versionType),
			zap.String("language", language),
			zap.Int64("prompt_id", po.ID))
		return po.ToEntity(), nil
	}

	if err != gorm.ErrRecordNotFound {
		logger.Error("failed to query style-generic prompt", zap.Error(err))
		return nil, fmt.Errorf("failed to query prompt: %w", err)
	}

	// 策略3：再降级（language='all' + style='all'）
	err = r.db.WithContext(ctx).
		Where("version_type = ? AND language = ? AND style = ? AND is_active = ?", versionType, "all", "all", true).
		Order("version DESC, weight DESC").
		First(&po).Error

	if err == nil {
		logger.Info("found language-generic prompt",
			zap.String("version_type", versionType),
			zap.Int64("prompt_id", po.ID))
		return po.ToEntity(), nil
	}

	if err == gorm.ErrRecordNotFound {
		logger.Warn("no active prompt found",
			zap.String("version_type", versionType),
			zap.String("language", language),
			zap.String("style", style))
		return nil, fmt.Errorf("no active prompt found for version_type=%s, language=%s, style=%s", versionType, language, style)
	}

	logger.Error("failed to query language-generic prompt", zap.Error(err))
	return nil, fmt.Errorf("failed to query prompt: %w", err)
}

// List 列出Prompts（支持过滤）
func (r *polishPromptRepositoryImpl) List(ctx context.Context, filter repository.PromptFilter) ([]*entity.PolishPrompt, error) {
	query := r.db.WithContext(ctx).Model(&PolishPromptPO{})

	// 应用过滤条件
	if filter.VersionType != "" {
		query = query.Where("version_type = ?", filter.VersionType)
	}
	if filter.Language != "" {
		query = query.Where("language = ?", filter.Language)
	}
	if filter.Style != "" {
		query = query.Where("style = ?", filter.Style)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.ABTestGroup != "" {
		query = query.Where("ab_test_group = ?", filter.ABTestGroup)
	}

	// 排序
	query = query.Order("version_type, language, style, version DESC")

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.Page > 1 {
			offset = (filter.Page - 1) * filter.PageSize
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	var pos []*PolishPromptPO
	if err := query.Find(&pos).Error; err != nil {
		logger.Error("failed to list polish prompts", zap.Error(err))
		return nil, fmt.Errorf("failed to list polish prompts: %w", err)
	}

	// 转换为实体
	prompts := make([]*entity.PolishPrompt, len(pos))
	for i, po := range pos {
		prompts[i] = po.ToEntity()
	}

	return prompts, nil
}

// Update 更新Prompt
func (r *polishPromptRepositoryImpl) Update(ctx context.Context, prompt *entity.PolishPrompt) error {
	po := &PolishPromptPO{}
	po.FromEntity(prompt)

	result := r.db.WithContext(ctx).Model(&PolishPromptPO{}).Where("id = ?", prompt.ID).Updates(po)
	if result.Error != nil {
		logger.Error("failed to update polish prompt", zap.Int64("id", prompt.ID), zap.Error(result.Error))
		return fmt.Errorf("failed to update polish prompt: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish prompt not found: id=%d", prompt.ID)
	}

	return nil
}

// Delete 删除Prompt（软删除，设置is_active=false）
func (r *polishPromptRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&PolishPromptPO{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		logger.Error("failed to delete polish prompt", zap.Int64("id", id), zap.Error(result.Error))
		return fmt.Errorf("failed to delete polish prompt: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish prompt not found: id=%d", id)
	}

	return nil
}

// Activate 激活Prompt
func (r *polishPromptRepositoryImpl) Activate(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&PolishPromptPO{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		logger.Error("failed to activate polish prompt", zap.Int64("id", id), zap.Error(result.Error))
		return fmt.Errorf("failed to activate polish prompt: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish prompt not found: id=%d", id)
	}

	return nil
}

// Deactivate 停用Prompt
func (r *polishPromptRepositoryImpl) Deactivate(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&PolishPromptPO{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		logger.Error("failed to deactivate polish prompt", zap.Int64("id", id), zap.Error(result.Error))
		return fmt.Errorf("failed to deactivate polish prompt: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish prompt not found: id=%d", id)
	}

	return nil
}

// IncrementUsage 增加使用次数
func (r *polishPromptRepositoryImpl) IncrementUsage(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&PolishPromptPO{}).Where("id = ?", id).UpdateColumn("usage_count", gorm.Expr("usage_count + 1"))
	if result.Error != nil {
		logger.Error("failed to increment prompt usage", zap.Int64("id", id), zap.Error(result.Error))
		return fmt.Errorf("failed to increment prompt usage: %w", result.Error)
	}

	return nil
}

// UpdateStatistics 更新统计信息
func (r *polishPromptRepositoryImpl) UpdateStatistics(ctx context.Context, id int64, successRate, avgSatisfaction float64) error {
	result := r.db.WithContext(ctx).Model(&PolishPromptPO{}).Where("id = ?", id).Updates(map[string]interface{}{
		"success_rate":     successRate,
		"avg_satisfaction": avgSatisfaction,
	})

	if result.Error != nil {
		logger.Error("failed to update prompt statistics", zap.Int64("id", id), zap.Error(result.Error))
		return fmt.Errorf("failed to update prompt statistics: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish prompt not found: id=%d", id)
	}

	return nil
}

// GetStatsByVersionType 按版本类型统计Prompt使用情况
func (r *polishPromptRepositoryImpl) GetStatsByVersionType(ctx context.Context) (map[string]*repository.PromptStats, error) {
	type statsResult struct {
		VersionType     string
		TotalCount      int64
		ActiveCount     int64
		TotalUsage      int64
		AvgSuccessRate  float64
		AvgSatisfaction float64
	}

	var results []statsResult
	err := r.db.WithContext(ctx).Model(&PolishPromptPO{}).
		Select(`
			version_type,
			COUNT(*) as total_count,
			SUM(CASE WHEN is_active = true THEN 1 ELSE 0 END) as active_count,
			SUM(usage_count) as total_usage,
			AVG(success_rate) as avg_success_rate,
			AVG(avg_satisfaction) as avg_satisfaction
		`).
		Group("version_type").
		Find(&results).Error

	if err != nil {
		logger.Error("failed to get prompt stats by version type", zap.Error(err))
		return nil, fmt.Errorf("failed to get prompt stats by version type: %w", err)
	}

	// 转换为map
	statsMap := make(map[string]*repository.PromptStats)
	for _, r := range results {
		statsMap[r.VersionType] = &repository.PromptStats{
			VersionType:     r.VersionType,
			TotalCount:      r.TotalCount,
			ActiveCount:     r.ActiveCount,
			TotalUsage:      r.TotalUsage,
			AvgSuccessRate:  r.AvgSuccessRate,
			AvgSatisfaction: r.AvgSatisfaction,
		}
	}

	return statsMap, nil
}
