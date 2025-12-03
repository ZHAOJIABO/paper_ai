package entity

import "time"

// PolishVersion 润色版本实体
// 用于存储多版本润色的每个版本详情
type PolishVersion struct {
	// 基础字段
	ID       int64
	RecordID int64 // 关联主表ID

	// 版本信息
	VersionType string // conservative / balanced / aggressive

	// 输出内容
	PolishedContent string
	PolishedLength  int
	Suggestions     []string // 改进建议

	// AI信息
	ModelUsed string
	PromptID  int64 // 关联使用的prompt模板ID

	// 性能指标
	ProcessTimeMs int

	// 状态
	Status       string // success / failed
	ErrorMessage string

	// 时间戳
	CreatedAt time.Time
}

// IsSuccess 判断是否成功
func (v *PolishVersion) IsSuccess() bool {
	return v.Status == "success"
}

// IsFailed 判断是否失败
func (v *PolishVersion) IsFailed() bool {
	return v.Status == "failed"
}

// GetChangeRate 获取修改率（与原文相比的长度变化百分比）
func (v *PolishVersion) GetChangeRate(originalLength int) float64 {
	if originalLength == 0 {
		return 0
	}
	return float64(v.PolishedLength-originalLength) / float64(originalLength) * 100
}

// VersionTypeEnum 版本类型枚举
const (
	VersionTypeConservative = "conservative" // 保守版本
	VersionTypeBalanced     = "balanced"     // 平衡版本
	VersionTypeAggressive   = "aggressive"   // 激进版本
)

// IsValidVersionType 验证版本类型是否有效
func IsValidVersionType(versionType string) bool {
	return versionType == VersionTypeConservative ||
		versionType == VersionTypeBalanced ||
		versionType == VersionTypeAggressive
}
