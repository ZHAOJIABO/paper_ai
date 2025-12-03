package repository

import (
	"context"
	"paper_ai/internal/domain/entity"
)

// PolishPromptRepository Prompt仓储接口
type PolishPromptRepository interface {
	// Create 创建Prompt
	Create(ctx context.Context, prompt *entity.PolishPrompt) error

	// GetByID 根据ID获取Prompt
	GetByID(ctx context.Context, id int64) (*entity.PolishPrompt, error)

	// GetActive 获取激活的Prompt (按版本类型、语言、风格查询)
	// 查询策略：
	// 1. 精确匹配：versionType + language + style
	// 2. 降级匹配：versionType + language + style='all'
	// 3. 再降级：versionType + language='all' + style='all'
	GetActive(ctx context.Context, versionType, language, style string) (*entity.PolishPrompt, error)

	// List 列出Prompts（支持过滤）
	List(ctx context.Context, filter PromptFilter) ([]*entity.PolishPrompt, error)

	// Update 更新Prompt
	Update(ctx context.Context, prompt *entity.PolishPrompt) error

	// Delete 删除Prompt（软删除）
	Delete(ctx context.Context, id int64) error

	// Activate 激活Prompt
	Activate(ctx context.Context, id int64) error

	// Deactivate 停用Prompt
	Deactivate(ctx context.Context, id int64) error

	// IncrementUsage 增加使用次数
	IncrementUsage(ctx context.Context, id int64) error

	// UpdateStatistics 更新统计信息
	UpdateStatistics(ctx context.Context, id int64, successRate, avgSatisfaction float64) error

	// GetStatsByVersionType 按版本类型统计Prompt使用情况
	GetStatsByVersionType(ctx context.Context) (map[string]*PromptStats, error)
}

// PromptFilter Prompt查询过滤器
type PromptFilter struct {
	VersionType string
	Language    string
	Style       string
	IsActive    *bool
	ABTestGroup string
	Page        int
	PageSize    int
}

// PromptStats Prompt统计
type PromptStats struct {
	VersionType     string
	TotalCount      int64
	ActiveCount     int64
	TotalUsage      int64
	AvgSuccessRate  float64
	AvgSatisfaction float64
}
