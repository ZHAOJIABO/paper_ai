package entity

import "time"

// PolishPrompt Prompt模板实体
// 用于数据库化管理Prompt模板，支持版本管理和A/B测试
type PolishPrompt struct {
	// 基础字段
	ID int64

	// 基本信息
	Name        string
	VersionType string // conservative / balanced / aggressive
	Language    string // en / zh / all
	Style       string // academic / formal / concise / all

	// Prompt内容
	SystemPrompt       string // 系统提示词
	UserPromptTemplate string // 用户提示词模板（支持变量替换）

	// 版本管理
	Version  int  // Prompt版本号
	IsActive bool // 是否启用

	// 元数据
	Description string
	Tags        []string // 标签

	// A/B测试
	ABTestGroup string // A/B测试分组
	Weight      int    // 权重（用于灰度发布）

	// 统计信息
	UsageCount      int     // 使用次数
	SuccessRate     float64 // 成功率
	AvgSatisfaction float64 // 平均满意度

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string // 创建人
}

// RenderUserPrompt 渲染用户提示词模板（替换变量）
func (p *PolishPrompt) RenderUserPrompt(variables map[string]string) string {
	prompt := p.UserPromptTemplate
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		// 简单的字符串替换，实际应用中可以使用更强大的模板引擎
		prompt = replaceAll(prompt, placeholder, value)
	}
	return prompt
}

// replaceAll 简单的字符串替换实现
func replaceAll(s, old, new string) string {
	// 使用 strings.ReplaceAll 的简单实现
	result := ""
	for {
		index := indexOf(s, old)
		if index == -1 {
			result += s
			break
		}
		result += s[:index] + new
		s = s[index+len(old):]
	}
	return result
}

// indexOf 查找子字符串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// IsEnabled 判断Prompt是否启用
func (p *PolishPrompt) IsEnabled() bool {
	return p.IsActive
}

// CanUseInABTest 判断是否可用于A/B测试
func (p *PolishPrompt) CanUseInABTest() bool {
	return p.ABTestGroup != "" && p.IsActive
}

// IncrementUsage 增加使用次数
func (p *PolishPrompt) IncrementUsage() {
	p.UsageCount++
}

// UpdateSuccessRate 更新成功率
func (p *PolishPrompt) UpdateSuccessRate(successCount, totalCount int) {
	if totalCount > 0 {
		p.SuccessRate = float64(successCount) / float64(totalCount) * 100
	}
}

// PromptLanguageEnum 语言枚举
const (
	PromptLanguageEnglish = "en"
	PromptLanguageChinese = "zh"
	PromptLanguageAll     = "all" // 通用
)

// PromptStyleEnum 风格枚举
const (
	PromptStyleAcademic = "academic"
	PromptStyleFormal   = "formal"
	PromptStyleConcise  = "concise"
	PromptStyleAll      = "all" // 通用
)

// IsValidLanguage 验证语言是否有效
func IsValidLanguage(language string) bool {
	return language == PromptLanguageEnglish ||
		language == PromptLanguageChinese ||
		language == PromptLanguageAll
}

// IsValidStyle 验证风格是否有效
func IsValidStyle(style string) bool {
	return style == PromptStyleAcademic ||
		style == PromptStyleFormal ||
		style == PromptStyleConcise ||
		style == PromptStyleAll
}
