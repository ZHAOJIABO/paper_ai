package types

// PolishRequest 润色请求
type PolishRequest struct {
	Content  string `json:"content"`  // 原始文本
	Style    string `json:"style"`    // 风格: academic/formal/concise
	Language string `json:"language"` // 语言: en/zh
}

// PolishResponse 润色响应
type PolishResponse struct {
	PolishedContent string   `json:"polished_content"` // 润色后的文本
	OriginalLength  int      `json:"original_length"`  // 原始长度
	PolishedLength  int      `json:"polished_length"`  // 润色后长度
	Suggestions     []string `json:"suggestions"`      // 改进建议
	ProviderUsed    string   `json:"provider_used"`    // 使用的提供商
	ModelUsed       string   `json:"model_used"`       // 使用的模型
}
