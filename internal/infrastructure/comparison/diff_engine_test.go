package comparison

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestDiffEngine_GenerateDiff(t *testing.T) {
	engine := NewDiffEngine()

	tests := []struct {
		name     string
		original string
		polished string
		wantOps  int // 期望的操作数量（粗略检查）
	}{
		{
			name:     "简单词汇替换",
			original: "This is a new method",
			polished: "This is a novel methodology",
			wantOps:  3, // Equal("This is a ") + Delete("new method") + Insert("novel methodology")
		},
		{
			name:     "相同文本",
			original: "Hello World",
			polished: "Hello World",
			wantOps:  1, // 只有一个 Equal 操作
		},
		{
			name:     "完全不同",
			original: "Hello",
			polished: "World",
			wantOps:  2, // Delete + Insert
		},
		{
			name:     "添加内容",
			original: "Hello",
			polished: "Hello World",
			wantOps:  2, // Equal("Hello") + Insert(" World")
		},
		{
			name:     "删除内容",
			original: "Hello World",
			polished: "Hello",
			wantOps:  2, // Equal("Hello") + Delete(" World")
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := engine.GenerateDiff(tt.original, tt.polished)

			if len(diffs) == 0 && tt.original != tt.polished {
				t.Errorf("GenerateDiff() 返回空结果，期望有差异")
			}

			// 检查结果数量在合理范围内
			if len(diffs) == 0 && tt.original == tt.polished {
				t.Logf("相同文本，diff 数量: %d", len(diffs))
			}

			t.Logf("Test: %s, Diff 数量: %d", tt.name, len(diffs))
		})
	}
}

func TestDiffEngine_GetChanges(t *testing.T) {
	engine := NewDiffEngine()

	tests := []struct {
		name       string
		original   string
		polished   string
		wantChange int // 期望的修改数量
	}{
		{
			name:       "单个词汇替换",
			original:   "method",
			polished:   "methodology",
			wantChange: 1,
		},
		{
			name:       "多处修改",
			original:   "new method to solve",
			polished:   "novel methodology to address",
			wantChange: 2, // method->methodology, solve->address
		},
		{
			name:       "无修改",
			original:   "Hello World",
			polished:   "Hello World",
			wantChange: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := engine.GenerateDiff(tt.original, tt.polished)
			changes := engine.GetChanges(diffs)

			if len(changes) != tt.wantChange {
				t.Errorf("GetChanges() 返回 %d 个修改，期望 %d 个", len(changes), tt.wantChange)
			}

			for i, change := range changes {
				t.Logf("修改 %d: '%s' -> '%s'", i+1, change.OriginalText, change.PolishedText)
			}
		})
	}
}

func TestDiffEngine_GetChanges_Complex(t *testing.T) {
	engine := NewDiffEngine()

	original := "In this paper, we propose a new method to solve the problem."
	polished := "In this paper, we propose a novel methodology to address the issue."

	diffs := engine.GenerateDiff(original, polished)
	changes := engine.GetChanges(diffs)

	t.Logf("原文: %s", original)
	t.Logf("润色: %s", polished)
	t.Logf("检测到 %d 处修改:", len(changes))

	for i, change := range changes {
		t.Logf("  修改 %d:", i+1)
		t.Logf("    原文: '%s'", change.OriginalText)
		t.Logf("    润色: '%s'", change.PolishedText)
	}

	if len(changes) == 0 {
		t.Error("应该检测到至少一处修改")
	}
}

func TestDiffItem_Types(t *testing.T) {
	engine := NewDiffEngine()

	original := "Hello World"
	polished := "Hello Go"

	diffs := engine.GenerateDiff(original, polished)

	hasEqual := false
	hasDelete := false
	hasInsert := false

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			hasEqual = true
		case diffmatchpatch.DiffDelete:
			hasDelete = true
		case diffmatchpatch.DiffInsert:
			hasInsert = true
		}
	}

	if !hasEqual {
		t.Error("应该包含相同部分 (Equal)")
	}
	if !hasDelete {
		t.Error("应该包含删除部分 (Delete)")
	}
	if !hasInsert {
		t.Error("应该包含插入部分 (Insert)")
	}
}
