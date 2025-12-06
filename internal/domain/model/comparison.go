package model

// ComparisonResult 对比结果
type ComparisonResult struct {
	TraceID         string      `json:"trace_id"`
	OriginalContent string      `json:"original_content"`
	PolishedContent string      `json:"polished_content"`
	FinalContent    string      `json:"final_content"`    // 用户应用修改后的最终内容
	Annotations     []Change    `json:"annotations"`
	Metadata        Metadata    `json:"metadata"`
	Statistics      Statistics  `json:"statistics"`
}

// Change 修改标注
type Change struct {
	ID               string        `json:"id"`                // 唯一标识
	Type             ChangeType    `json:"type"`              // vocabulary/grammar/structure

	// 润色后文本中的位置（前端高亮用）
	PolishedPosition Position      `json:"polished_position"`
	PolishedText     string        `json:"polished_text"`     // 修改后的文本

	// 原文信息（悬浮时展示）
	OriginalText     string        `json:"original_text"`     // 原始文本

	// 详情信息（右侧面板展示）
	Reason           string        `json:"reason"`            // 修改理由
	Alternatives     []Alternative `json:"alternatives"`      // 替代方案
	Confidence       float64       `json:"confidence"`        // 置信度 0-1
	Impact           string        `json:"impact"`            // 影响维度
	HighlightColor   string        `json:"highlight_color"`   // 建议的高亮颜色

	// 用户操作状态
	Status           ActionStatus  `json:"status"`            // pending/accepted/rejected
}

// Position 位置信息
type Position struct {
	Start int `json:"start"` // 字符级起始位置（基于 rune）
	End   int `json:"end"`   // 字符级结束位置
	Line  int `json:"line"`  // 所在行号（可选）
}

// Alternative 替代方案
type Alternative struct {
	Text   string `json:"text"`   // 替代文本
	Reason string `json:"reason"` // 使用理由
}

// Metadata 元数据
type Metadata struct {
	OriginalWordCount        int     `json:"original_word_count"`
	PolishedWordCount        int     `json:"polished_word_count"`
	TotalChanges             int     `json:"total_changes"`
	AcademicScoreImprovement float64 `json:"academic_score_improvement"` // 百分比
}

// Statistics 统计信息
type Statistics struct {
	VocabularyChanges int `json:"vocabulary_changes"`
	GrammarChanges    int `json:"grammar_changes"`
	StructureChanges  int `json:"structure_changes"`
}

// ChangeType 修改类型
type ChangeType string

const (
	ChangeTypeVocabulary ChangeType = "vocabulary" // 词汇优化
	ChangeTypeGrammar    ChangeType = "grammar"    // 语法修正
	ChangeTypeStructure  ChangeType = "structure"  // 结构调整
)

// ActionStatus 操作状态
type ActionStatus string

const (
	ActionStatusPending  ActionStatus = "pending"  // 待处理
	ActionStatusAccepted ActionStatus = "accepted" // 已接受
	ActionStatusRejected ActionStatus = "rejected" // 已拒绝
)

// ChangeActionRequest 修改操作请求
type ChangeActionRequest struct {
	ChangeID         string `json:"change_id" binding:"required"`
	Action           string `json:"action" binding:"required,oneof=accept reject"`
	AlternativeIndex *int   `json:"alternative_index"` // 可选：选择替代方案
}

// ChangeActionResponse 修改操作响应
type ChangeActionResponse struct {
	Success        bool     `json:"success"`
	UpdatedContent string   `json:"updated_content"`
	AppliedChanges []string `json:"applied_changes"`
	PendingChanges []string `json:"pending_changes"`
}

// BatchActionRequest 批量操作请求
type BatchActionRequest struct {
	Action    string   `json:"action" binding:"required,oneof=accept_all reject_all"`
	ChangeIDs []string `json:"change_ids"` // 可选：指定特定的修改
}

// BatchActionResponse 批量操作响应
type BatchActionResponse struct {
	Success      bool   `json:"success"`
	UpdatedContent string `json:"updated_content"`
	AppliedCount int    `json:"applied_count"`
}
