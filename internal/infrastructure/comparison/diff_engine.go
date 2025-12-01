package comparison

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffItem 差异项
type DiffItem struct {
	Type diffmatchpatch.Operation // 操作类型：DiffDelete, DiffInsert, DiffEqual
	Text string                    // 文本内容
}

// DiffEngine 差异引擎
type DiffEngine struct {
	dmp *diffmatchpatch.DiffMatchPatch
}

// NewDiffEngine 创建差异引擎
func NewDiffEngine() *DiffEngine {
	return &DiffEngine{
		dmp: diffmatchpatch.New(),
	}
}

// GenerateDiff 生成文本差异
func (e *DiffEngine) GenerateDiff(original, polished string) []DiffItem {
	// 1. 运行 diff 算法
	diffs := e.dmp.DiffMain(original, polished, false)

	// 2. 优化 diff 结果（合并语义相关的改动）
	diffs = e.dmp.DiffCleanupSemantic(diffs)

	// 3. 转换为内部格式
	items := make([]DiffItem, 0, len(diffs))
	for _, diff := range diffs {
		items = append(items, DiffItem{
			Type: diff.Type,
			Text: diff.Text,
		})
	}

	return items
}

// GetChanges 从 diff 结果中提取修改对
func (e *DiffEngine) GetChanges(diffs []DiffItem) []ChangeInfo {
	changes := make([]ChangeInfo, 0)

	i := 0
	for i < len(diffs) {
		// 查找删除-插入对（表示替换）
		if i < len(diffs)-1 {
			if diffs[i].Type == diffmatchpatch.DiffDelete && diffs[i+1].Type == diffmatchpatch.DiffInsert {
				changes = append(changes, ChangeInfo{
					OriginalText: diffs[i].Text,
					PolishedText: diffs[i+1].Text,
				})
				i += 2
				continue
			}
		}

		// 单独的删除
		if diffs[i].Type == diffmatchpatch.DiffDelete {
			changes = append(changes, ChangeInfo{
				OriginalText: diffs[i].Text,
				PolishedText: "",
			})
		}

		// 单独的插入
		if diffs[i].Type == diffmatchpatch.DiffInsert {
			changes = append(changes, ChangeInfo{
				OriginalText: "",
				PolishedText: diffs[i].Text,
			})
		}

		i++
	}

	return changes
}

// ChangeInfo 修改信息
type ChangeInfo struct {
	OriginalText string
	PolishedText string
}
