package model

// PolishRequest 段落润色请求模型
type PolishRequest struct {
	Content  string `json:"content" binding:"required"`
	Provider string `json:"provider"`
	Style    string `json:"style"`
	Language string `json:"language"`
}

// Validate 验证请求参数
func (r *PolishRequest) Validate() error {
	if r.Content == "" {
		return &ValidationError{Field: "content", Message: "content cannot be empty"}
	}

	if len(r.Content) > 10000 {
		return &ValidationError{Field: "content", Message: "content too long, maximum 10000 characters"}
	}

	// 验证style
	if r.Style != "" && !isValidStyle(r.Style) {
		return &ValidationError{Field: "style", Message: "invalid style, must be one of: academic, formal, concise"}
	}

	// 验证language
	if r.Language != "" && !isValidLanguage(r.Language) {
		return &ValidationError{Field: "language", Message: "invalid language, must be one of: en, zh"}
	}

	return nil
}

// SetDefaults 设置默认值
func (r *PolishRequest) SetDefaults() {
	if r.Style == "" {
		r.Style = "academic"
	}
	if r.Language == "" {
		r.Language = "en"
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// isValidStyle 验证风格是否有效
func isValidStyle(style string) bool {
	validStyles := map[string]bool{
		"academic": true,
		"formal":   true,
		"concise":  true,
	}
	return validStyles[style]
}

// isValidLanguage 验证语言是否有效
func isValidLanguage(lang string) bool {
	validLanguages := map[string]bool{
		"en": true,
		"zh": true,
	}
	return validLanguages[lang]
}
