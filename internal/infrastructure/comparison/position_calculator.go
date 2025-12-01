package comparison

import (
	"strings"
)

// PositionCalculator 位置计算器
type PositionCalculator struct{}

// NewPositionCalculator 创建位置计算器
func NewPositionCalculator() *PositionCalculator {
	return &PositionCalculator{}
}

// CalculatePositions 计算所有修改在润色后文本中的位置
func (c *PositionCalculator) CalculatePositions(polishedText string, changes []ChangeInfo) []PositionInfo {
	positions := make([]PositionInfo, 0, len(changes))
	runes := []rune(polishedText)
	currentPos := 0

	for _, change := range changes {
		// 如果是删除（只在原文中存在），跳过
		if change.PolishedText == "" {
			continue
		}

		// 在润色后文本中查找修改的位置
		targetRunes := []rune(change.PolishedText)
		pos := c.findPosition(runes, targetRunes, currentPos)

		if pos.Start >= 0 {
			positions = append(positions, PositionInfo{
				Start:        pos.Start,
				End:          pos.End,
				Line:         pos.Line,
				OriginalText: change.OriginalText,
				PolishedText: change.PolishedText,
			})
			currentPos = pos.End // 更新搜索起始位置
		}
	}

	return positions
}

// findPosition 在文本中查找目标文本的位置
func (c *PositionCalculator) findPosition(runes, targetRunes []rune, startPos int) Position {
	for i := startPos; i <= len(runes)-len(targetRunes); i++ {
		match := true
		for j := 0; j < len(targetRunes); j++ {
			if runes[i+j] != targetRunes[j] {
				match = false
				break
			}
		}

		if match {
			return Position{
				Start: i,
				End:   i + len(targetRunes),
				Line:  c.calculateLineNumber(runes, i),
			}
		}
	}

	return Position{Start: -1, End: -1, Line: -1}
}

// calculateLineNumber 计算位置所在的行号
func (c *PositionCalculator) calculateLineNumber(runes []rune, position int) int {
	if position < 0 || position >= len(runes) {
		return -1
	}

	lineCount := 1
	for i := 0; i < position && i < len(runes); i++ {
		if runes[i] == '\n' {
			lineCount++
		}
	}
	return lineCount
}

// Position 位置信息（内部使用）
type Position struct {
	Start int
	End   int
	Line  int
}

// PositionInfo 带文本的位置信息
type PositionInfo struct {
	Start        int
	End          int
	Line         int
	OriginalText string
	PolishedText string
}

// CountWords 统计文本中的词数
func CountWords(text string) int {
	// 简单实现：按空格分割
	words := strings.Fields(text)
	return len(words)
}
