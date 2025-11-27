package entity

import "time"

// PolishRecord 润色记录实体
// 纯业务模型，不包含任何ORM框架标签，保持领域层的纯净性
type PolishRecord struct {
	// 基础字段
	ID      int64
	TraceID string

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

	// 性能指标
	ProcessTimeMs int

	// 状态信息
	Status       string
	ErrorMessage string

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
