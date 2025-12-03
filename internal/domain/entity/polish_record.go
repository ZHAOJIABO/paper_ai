package entity

import "time"

// PolishRecord 润色记录实体
// 纯业务模型，不包含任何ORM框架标签，保持领域层的纯净性
type PolishRecord struct {
	// 基础字段
	ID      int64
	TraceID string
	UserID  int64 // 用户ID

	// 输入信息
	OriginalContent string
	Style           string
	Language        string

	// 输出信息
	PolishedContent string
	OriginalLength  int
	PolishedLength  int

	// AI信息
	Provider string
	Model    string

	// 模式信息
	Mode string // single / multi

	// 性能指标
	ProcessTimeMs int

	// 状态信息
	Status       string
	ErrorMessage string

	// 对比数据（新增）
	ComparisonData  string   // 对比数据JSON
	ChangesCount    int      // 修改总数
	AcceptedChanges []string // 用户接受的修改ID列表
	RejectedChanges []string // 用户拒绝的修改ID列表
	FinalContent    string   // 用户最终确认的文本（原文+接受的修改）

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsSuccess 判断是否成功
func (r *PolishRecord) IsSuccess() bool {
	return r.Status == "success"
}

// IsFailed 判断是否失败
func (r *PolishRecord) IsFailed() bool {
	return r.Status == "failed"
}

// GetContentDiff 获取内容长度差异
func (r *PolishRecord) GetContentDiff() int {
	return r.PolishedLength - r.OriginalLength
}

// GetCompressionRate 获取压缩率（负数表示内容变长）
func (r *PolishRecord) GetCompressionRate() float64 {
	if r.OriginalLength == 0 {
		return 0
	}
	return float64(r.OriginalLength-r.PolishedLength) / float64(r.OriginalLength) * 100
}

// IsMultiVersionMode 判断是否为多版本模式
func (r *PolishRecord) IsMultiVersionMode() bool {
	return r.Mode == "multi"
}

// IsSingleVersionMode 判断是否为单版本模式
func (r *PolishRecord) IsSingleVersionMode() bool {
	return r.Mode == "single" || r.Mode == ""
}

// ModeEnum 模式枚举
const (
	ModeSingle = "single" // 单版本模式
	ModeMulti  = "multi"  // 多版本模式
)

// IsValidMode 验证模式是否有效
func IsValidMode(mode string) bool {
	return mode == ModeSingle || mode == ModeMulti
}
