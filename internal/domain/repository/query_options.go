package repository

import "time"

// QueryOptions 查询选项（Options模式，高扩展性）
type QueryOptions struct {
	// 分页
	Page     int
	PageSize int
	Offset   int
	Limit    int

	// 过滤条件
	UserID   *int64  // 按用户ID过滤 ⭐ 新增
	Provider *string // 按提供商过滤
	Status   *string // 按状态过滤
	Language *string // 按语言过滤
	Style    *string // 按风格过滤

	// 时间范围
	StartTime *time.Time
	EndTime   *time.Time

	// 排序
	OrderBy   string // 排序字段
	OrderDesc bool   // 是否降序

	// 字段选择（性能优化）
	SelectFields []string // 需要返回的字段
	ExcludeText  bool     // 是否排除大文本字段（original_content, polished_content）
}

// QueryOptionsBuilder 查询选项构建器（Fluent API）
type QueryOptionsBuilder struct {
	opts QueryOptions
}

// NewQueryOptions 创建查询选项构建器
func NewQueryOptions() *QueryOptionsBuilder {
	return &QueryOptionsBuilder{
		opts: QueryOptions{
			OrderBy:   "created_at",
			OrderDesc: true, // 默认按创建时间降序
		},
	}
}

// Page 设置分页
func (b *QueryOptionsBuilder) Page(page, pageSize int) *QueryOptionsBuilder {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页大小
	}

	b.opts.Page = page
	b.opts.PageSize = pageSize
	b.opts.Offset = (page - 1) * pageSize
	b.opts.Limit = pageSize
	return b
}

// WithUserID 按用户ID过滤 ⭐ 新增
func (b *QueryOptionsBuilder) WithUserID(userID int64) *QueryOptionsBuilder {
	b.opts.UserID = &userID
	return b
}

// WithProvider 按提供商过滤
func (b *QueryOptionsBuilder) WithProvider(provider string) *QueryOptionsBuilder {
	b.opts.Provider = &provider
	return b
}

// WithStatus 按状态过滤
func (b *QueryOptionsBuilder) WithStatus(status string) *QueryOptionsBuilder {
	b.opts.Status = &status
	return b
}

// WithLanguage 按语言过滤
func (b *QueryOptionsBuilder) WithLanguage(language string) *QueryOptionsBuilder {
	b.opts.Language = &language
	return b
}

// WithStyle 按风格过滤
func (b *QueryOptionsBuilder) WithStyle(style string) *QueryOptionsBuilder {
	b.opts.Style = &style
	return b
}

// WithTimeRange 时间范围过滤
func (b *QueryOptionsBuilder) WithTimeRange(start, end time.Time) *QueryOptionsBuilder {
	b.opts.StartTime = &start
	b.opts.EndTime = &end
	return b
}

// OrderBy 排序
func (b *QueryOptionsBuilder) OrderBy(field string, desc bool) *QueryOptionsBuilder {
	b.opts.OrderBy = field
	b.opts.OrderDesc = desc
	return b
}

// SelectFields 选择返回字段
func (b *QueryOptionsBuilder) SelectFields(fields []string) *QueryOptionsBuilder {
	b.opts.SelectFields = fields
	return b
}

// ExcludeText 排除大文本字段（性能优化）
func (b *QueryOptionsBuilder) ExcludeText() *QueryOptionsBuilder {
	b.opts.ExcludeText = true
	return b
}

// Build 构建查询选项
func (b *QueryOptionsBuilder) Build() QueryOptions {
	return b.opts
}

// StatisticsOptions 统计选项
type StatisticsOptions struct {
	UserID      *int64     // 按用户ID过滤
	TimeRange   *TimeRange
	GroupBy     []string // 分组维度：provider, status, language, style
	Aggregation []string // 聚合字段：count, avg_time, success_rate
}

// TimeRange 时间范围
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Statistics 统计结果
type Statistics struct {
	TotalCount     int64            `json:"total_count"`
	SuccessCount   int64            `json:"success_count"`
	FailedCount    int64            `json:"failed_count"`
	SuccessRate    float64          `json:"success_rate"`
	AvgProcessTime float64          `json:"avg_process_time_ms"`
	ProviderStats  map[string]*Stat `json:"provider_stats,omitempty"`
	LanguageStats  map[string]*Stat `json:"language_stats,omitempty"`
	StyleStats     map[string]*Stat `json:"style_stats,omitempty"`
}

// Stat 单项统计
type Stat struct {
	Count          int64   `json:"count"`
	SuccessCount   int64   `json:"success_count"`
	FailedCount    int64   `json:"failed_count"`
	SuccessRate    float64 `json:"success_rate"`
	AvgProcessTime float64 `json:"avg_process_time_ms"`
}
