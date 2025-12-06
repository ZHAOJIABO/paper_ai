package model

// PolishMultiVersionRequest 多版本润色请求
type PolishMultiVersionRequest struct {
	Content  string   `json:"content"`  // 原始内容
	Style    string   `json:"style"`    // 润色风格: academic/formal/concise
	Language string   `json:"language"` // 语言: en/zh
	Provider string   `json:"provider"` // AI提供商: claude/doubao等
	Versions []string `json:"versions"` // 指定需要的版本类型，不指定则生成全部3个版本
}

// PolishMultiVersionResponse 多版本润色响应
type PolishMultiVersionResponse struct {
	TraceID         string                    `json:"trace_id"`         // 追踪ID
	OriginalContent string                    `json:"original_content"` // 原始内容
	OriginalLength  int                       `json:"original_length"`  // 原文长度
	Versions        map[string]*VersionResult `json:"versions"`         // 版本结果: version_type -> result
	ProviderUsed    string                    `json:"provider_used"`    // 使用的提供商
}

// VersionResult 单个版本的润色结果
type VersionResult struct {
	PolishedContent string   `json:"polished_content"` // 润色后的内容
	PolishedLength  int      `json:"polished_length"`  // 润色后的长度
	Suggestions     []string `json:"suggestions"`      // 改进建议
	ProcessTimeMs   int      `json:"process_time_ms"`  // 处理耗时(毫秒)
	ModelUsed       string   `json:"model_used"`       // 使用的模型
	Status          string   `json:"status"`           // 状态: success / failed
	ErrorMessage    string   `json:"error_message"`    // 错误信息(如果失败)
}
