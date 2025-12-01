package comparison

import (
	"testing"
)

func TestPositionCalculator_CalculatePositions(t *testing.T) {
	calc := NewPositionCalculator()

	tests := []struct {
		name          string
		polishedText  string
		changes       []ChangeInfo
		wantPositions int
	}{
		{
			name:         "单个修改",
			polishedText: "This is a novel methodology",
			changes: []ChangeInfo{
				{OriginalText: "method", PolishedText: "methodology"},
			},
			wantPositions: 1,
		},
		{
			name:         "多个修改",
			polishedText: "In this paper, we propose a novel methodology to address the issue.",
			changes: []ChangeInfo{
				{OriginalText: "method", PolishedText: "methodology"},
				{OriginalText: "solve", PolishedText: "address"},
			},
			wantPositions: 2,
		},
		{
			name:         "删除操作（润色文本中不存在）",
			polishedText: "Hello World",
			changes: []ChangeInfo{
				{OriginalText: "something", PolishedText: ""},
			},
			wantPositions: 0, // 删除操作不应该在润色文本中出现
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positions := calc.CalculatePositions(tt.polishedText, tt.changes)

			if len(positions) != tt.wantPositions {
				t.Errorf("CalculatePositions() 返回 %d 个位置，期望 %d 个", len(positions), tt.wantPositions)
			}

			for i, pos := range positions {
				t.Logf("位置 %d:", i+1)
				t.Logf("  Start: %d, End: %d, Line: %d", pos.Start, pos.End, pos.Line)
				t.Logf("  原文: '%s'", pos.OriginalText)
				t.Logf("  润色: '%s'", pos.PolishedText)

				// 验证位置范围
				if pos.Start < 0 || pos.End > len([]rune(tt.polishedText)) {
					t.Errorf("位置超出范围: Start=%d, End=%d, TextLen=%d", pos.Start, pos.End, len([]rune(tt.polishedText)))
				}

				// 验证提取的文本
				runes := []rune(tt.polishedText)
				extractedText := string(runes[pos.Start:pos.End])
				if extractedText != pos.PolishedText {
					t.Errorf("位置提取错误: 期望 '%s', 得到 '%s'", pos.PolishedText, extractedText)
				}
			}
		})
	}
}

func TestPositionCalculator_LineNumber(t *testing.T) {
	calc := NewPositionCalculator()

	polishedText := "Line 1: Hello\nLine 2: World\nLine 3: Test"
	changes := []ChangeInfo{
		{OriginalText: "Hi", PolishedText: "Hello"},
		{OriginalText: "Earth", PolishedText: "World"},
		{OriginalText: "Check", PolishedText: "Test"},
	}

	positions := calc.CalculatePositions(polishedText, changes)

	expectedLines := []int{1, 2, 3}

	for i, pos := range positions {
		if i >= len(expectedLines) {
			break
		}
		if pos.Line != expectedLines[i] {
			t.Errorf("修改 %d 行号错误: 期望 %d, 得到 %d", i+1, expectedLines[i], pos.Line)
		}
	}
}

func TestPositionCalculator_Unicode(t *testing.T) {
	calc := NewPositionCalculator()

	// 测试 Unicode 字符（中文、emoji）
	polishedText := "这是一个测试文本 with emoji 😀"
	changes := []ChangeInfo{
		{OriginalText: "测试", PolishedText: "测试"},
		{OriginalText: "😀", PolishedText: "😀"},
	}

	positions := calc.CalculatePositions(polishedText, changes)

	if len(positions) != 2 {
		t.Errorf("Unicode 文本位置计算错误: 期望 2 个位置，得到 %d 个", len(positions))
	}

	for i, pos := range positions {
		runes := []rune(polishedText)
		extractedText := string(runes[pos.Start:pos.End])
		if extractedText != pos.PolishedText {
			t.Errorf("Unicode 位置 %d 提取错误: 期望 '%s', 得到 '%s'", i+1, pos.PolishedText, extractedText)
		}
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantCount int
	}{
		{
			name:      "简单句子",
			text:      "Hello World",
			wantCount: 2,
		},
		{
			name:      "多个空格",
			text:      "Hello   World  Test",
			wantCount: 3,
		},
		{
			name:      "空字符串",
			text:      "",
			wantCount: 0,
		},
		{
			name:      "只有空格",
			text:      "   ",
			wantCount: 0,
		},
		{
			name:      "包含标点",
			text:      "Hello, World! How are you?",
			wantCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := CountWords(tt.text)
			if count != tt.wantCount {
				t.Errorf("CountWords() = %d, 期望 %d", count, tt.wantCount)
			}
		})
	}
}

func TestPositionCalculator_MultipleOccurrences(t *testing.T) {
	calc := NewPositionCalculator()

	// 测试同一个词出现多次的情况
	polishedText := "test test test"
	changes := []ChangeInfo{
		{OriginalText: "check", PolishedText: "test"},
	}

	positions := calc.CalculatePositions(polishedText, changes)

	// 应该只找到第一个出现的位置
	if len(positions) != 1 {
		t.Errorf("应该只找到一个位置，得到 %d 个", len(positions))
	}

	if len(positions) > 0 {
		if positions[0].Start != 0 {
			t.Errorf("应该找到第一个出现的位置(0)，得到 %d", positions[0].Start)
		}
	}
}
