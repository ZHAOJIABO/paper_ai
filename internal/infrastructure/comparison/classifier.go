package comparison

import (
	"math"
	"strings"
	"paper_ai/internal/domain/model"
)

// ChangeClassifier 修改分类器
type ChangeClassifier struct{}

// NewChangeClassifier 创建修改分类器
func NewChangeClassifier() *ChangeClassifier {
	return &ChangeClassifier{}
}

// Classify 对修改进行分类
func (c *ChangeClassifier) Classify(originalText, polishedText string) model.ChangeType {
	// 处理空文本情况
	if originalText == "" || polishedText == "" {
		return model.ChangeTypeStructure
	}

	// 规则3: 结构调整（词数变化超过30%）
	if c.isStructuralChange(originalText, polishedText) {
		return model.ChangeTypeStructure
	}

	// 规则2: 语法修正
	if c.isGrammarFix(originalText, polishedText) {
		return model.ChangeTypeGrammar
	}

	// 规则1: 词汇优化（默认）
	return model.ChangeTypeVocabulary
}

// isSynonymReplacement 检测是否为同义词替换
func (c *ChangeClassifier) isSynonymReplacement(orig, pol string) bool {
	origWords := strings.Fields(orig)
	polWords := strings.Fields(pol)

	// 词数相同，且都是短语（1-3个词）
	if len(origWords) == len(polWords) && len(origWords) <= 3 {
		return true
	}
	return false
}

// isGrammarFix 检测是否为语法修正
func (c *ChangeClassifier) isGrammarFix(orig, pol string) bool {
	// 语法关键词列表（冠词、介词、时态助动词等）
	grammarKeywords := []string{
		"a", "an", "the",           // 冠词
		"is", "are", "was", "were", // be 动词
		"in", "on", "at", "to", "for", "with", // 介词
		"have", "has", "had",       // have 动词
	}

	origLower := strings.ToLower(orig)
	polLower := strings.ToLower(pol)

	// 检查是否包含语法关键词
	for _, keyword := range grammarKeywords {
		if strings.Contains(origLower, " "+keyword+" ") || strings.Contains(polLower, " "+keyword+" ") {
			return true
		}
		// 检查开头和结尾
		if strings.HasPrefix(origLower, keyword+" ") || strings.HasPrefix(polLower, keyword+" ") {
			return true
		}
		if strings.HasSuffix(origLower, " "+keyword) || strings.HasSuffix(polLower, " "+keyword) {
			return true
		}
	}

	return false
}

// isStructuralChange 检测是否为结构调整
func (c *ChangeClassifier) isStructuralChange(orig, pol string) bool {
	origWords := len(strings.Fields(orig))
	polWords := len(strings.Fields(pol))

	// 避免除以零
	if origWords == 0 {
		return polWords > 0
	}

	// 词数变化超过30%，判定为结构调整
	diff := math.Abs(float64(origWords - polWords))
	ratio := diff / float64(origWords)

	return ratio > 0.3
}

// SuggestHighlightColor 根据修改类型建议高亮颜色
func (c *ChangeClassifier) SuggestHighlightColor(changeType model.ChangeType) string {
	switch changeType {
	case model.ChangeTypeVocabulary:
		return "yellow" // 黄色：词汇优化
	case model.ChangeTypeGrammar:
		return "lightblue" // 蓝色：语法修正
	case model.ChangeTypeStructure:
		return "lightgreen" // 绿色：结构调整
	default:
		return "yellow"
	}
}
