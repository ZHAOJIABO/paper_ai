package model

// PolishMultiVersionRequest 多版本润色请求
type PolishMultiVersionRequest struct {
	Content  string   // 原始内容
	Style    string   // 润色风格: academic/formal/concise
	Language string   // 语言: en/zh
	Provider string   // AI提供商: claude/doubao等
	Versions []string // 指定需要的版本类型，不指定则生成全部3个版本
}

// PolishMultiVersionResponse 多版本润色响应
type PolishMultiVersionResponse struct {
	TraceID        string                   // 追踪ID
	OriginalLength int                      // 原文长度
	Versions       map[string]*VersionResult // 版本结果: version_type -> result
	ProviderUsed   string                   // 使用的提供商
}

// VersionResult 单个版本的润色结果
type VersionResult struct {
	PolishedContent string   // 润色后的内容
	PolishedLength  int      // 润色后的长度
	Suggestions     []string // 改进建议
	ProcessTimeMs   int      // 处理耗时(毫秒)
	ModelUsed       string   // 使用的模型
	Status          string   // 状态: success / failed
	ErrorMessage    string   // 错误信息(如果失败)
}
