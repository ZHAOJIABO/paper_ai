package comparison

import (
	"fmt"
	"paper_ai/internal/domain/model"
)

// ReasonGenerator 理由生成器
type ReasonGenerator struct{}

// NewReasonGenerator 创建理由生成器
func NewReasonGenerator() *ReasonGenerator {
	return &ReasonGenerator{}
}

// Generate 生成修改理由
func (g *ReasonGenerator) Generate(changeType model.ChangeType, original, polished string) string {
	switch changeType {
	case model.ChangeTypeVocabulary:
		if original != "" && polished != "" {
			return fmt.Sprintf("使用更学术化的词汇，在学术写作中 '%s' 比 '%s' 更正式和专业，能更准确地表达研究内容",
				polished, original)
		}
		return "优化了词汇选择，提升了学术表达的准确性和专业性"

	case model.ChangeTypeGrammar:
		return "语法优化，修正了表达方式，使句子更符合学术写作规范"

	case model.ChangeTypeStructure:
		return "结构调整，重新组织句子结构，提升了可读性和逻辑连贯性"

	default:
		return "优化了表达方式，提升了文本质量"
	}
}

// GenerateAlternatives 生成替代方案
func (g *ReasonGenerator) GenerateAlternatives(changeType model.ChangeType, original string) []model.Alternative {
	// 简单实现：返回预定义的替代方案
	// TODO: 未来可以调用 AI 生成更智能的替代方案

	alternatives := []model.Alternative{}

	switch changeType {
	case model.ChangeTypeVocabulary:
		// 根据常见词汇提供替代方案
		alternatives = g.getVocabularyAlternatives(original)

	case model.ChangeTypeGrammar:
		// 语法修正通常不需要替代方案
		alternatives = []model.Alternative{}

	case model.ChangeTypeStructure:
		// 结构调整较复杂，暂不提供替代方案
		alternatives = []model.Alternative{}
	}

	return alternatives
}

// getVocabularyAlternatives 获取词汇类替代方案
func (g *ReasonGenerator) getVocabularyAlternatives(original string) []model.Alternative {
	// 预定义的常见学术词汇映射
	vocabularyMap := map[string][]model.Alternative{
		"method": {
			{Text: "approach", Reason: "更通用的学术表达，适用于描述研究路径"},
			{Text: "technique", Reason: "强调技术性方法，适合技术类论文"},
		},
		"problem": {
			{Text: "issue", Reason: "更正式的表达"},
			{Text: "challenge", Reason: "强调挑战性"},
		},
		"solve": {
			{Text: "address", Reason: "更学术化的表达"},
			{Text: "tackle", Reason: "更积极主动的表达"},
		},
		"show": {
			{Text: "demonstrate", Reason: "更正式的学术表达"},
			{Text: "illustrate", Reason: "强调清晰展示"},
		},
		"use": {
			{Text: "utilize", Reason: "更正式的学术用语"},
			{Text: "employ", Reason: "强调有意识地使用"},
		},
		"get": {
			{Text: "obtain", Reason: "更正式的表达"},
			{Text: "acquire", Reason: "强调获取过程"},
		},
		"make": {
			{Text: "construct", Reason: "强调构建过程"},
			{Text: "develop", Reason: "强调开发过程"},
		},
		"help": {
			{Text: "facilitate", Reason: "更学术化的表达"},
			{Text: "assist", Reason: "更正式的表达"},
		},
	}

	// 查找是否有匹配的替代方案
	if alternatives, ok := vocabularyMap[original]; ok {
		return alternatives
	}

	// 如果没有预定义的，返回通用建议
	return []model.Alternative{
		{Text: "approach", Reason: "通用的学术表达"},
	}
}

// CalculateConfidence 计算修改的置信度
func (g *ReasonGenerator) CalculateConfidence(changeType model.ChangeType, original, polished string) float64 {
	// 简单实现：根据修改类型返回不同的置信度
	switch changeType {
	case model.ChangeTypeVocabulary:
		// 词汇替换置信度较高
		return 0.90

	case model.ChangeTypeGrammar:
		// 语法修正置信度中等
		return 0.85

	case model.ChangeTypeStructure:
		// 结构调整置信度相对较低
		return 0.75

	default:
		return 0.80
	}
}

// GetImpact 获取修改的影响维度
func (g *ReasonGenerator) GetImpact(changeType model.ChangeType) string {
	switch changeType {
	case model.ChangeTypeVocabulary:
		return "academic_tone" // 学术语气

	case model.ChangeTypeGrammar:
		return "grammar_correctness" // 语法正确性

	case model.ChangeTypeStructure:
		return "readability" // 可读性

	default:
		return "general" // 通用
	}
}
