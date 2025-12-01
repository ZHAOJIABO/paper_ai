package comparison

import (
	"testing"

	"paper_ai/internal/domain/model"
)

func TestChangeClassifier_Classify(t *testing.T) {
	classifier := NewChangeClassifier()

	tests := []struct {
		name         string
		originalText string
		polishedText string
		wantType     model.ChangeType
	}{
		{
			name:         "词汇优化 - 同义词替换",
			originalText: "method",
			polishedText: "methodology",
			wantType:     model.ChangeTypeVocabulary,
		},
		{
			name:         "词汇优化 - 短语替换",
			originalText: "new method",
			polishedText: "novel approach",
			wantType:     model.ChangeTypeVocabulary,
		},
		{
			name:         "语法修正 - 冠词",
			originalText: "a apple",
			polishedText: "an apple",
			wantType:     model.ChangeTypeGrammar,
		},
		{
			name:         "语法修正 - be动词",
			originalText: "he are",
			polishedText: "he is",
			wantType:     model.ChangeTypeGrammar,
		},
		{
			name:         "结构调整 - 词数大幅变化",
			originalText: "hello",
			polishedText: "hello world this is a test",
			wantType:     model.ChangeTypeStructure,
		},
		{
			name:         "空文本 - 删除",
			originalText: "something",
			polishedText: "",
			wantType:     model.ChangeTypeStructure,
		},
		{
			name:         "空文本 - 添加",
			originalText: "",
			polishedText: "something new",
			wantType:     model.ChangeTypeStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType := classifier.Classify(tt.originalText, tt.polishedText)

			if gotType != tt.wantType {
				t.Errorf("Classify() = %v, 期望 %v", gotType, tt.wantType)
			}

			t.Logf("原文: '%s' -> 润色: '%s' => 类型: %s", tt.originalText, tt.polishedText, gotType)
		})
	}
}

func TestChangeClassifier_isSynonymReplacement(t *testing.T) {
	classifier := NewChangeClassifier()

	tests := []struct {
		name     string
		original string
		polished string
		want     bool
	}{
		{
			name:     "单词替换",
			original: "method",
			polished: "approach",
			want:     true,
		},
		{
			name:     "两个词替换",
			original: "new method",
			polished: "novel approach",
			want:     true,
		},
		{
			name:     "词数不同",
			original: "method",
			polished: "new approach",
			want:     false,
		},
		{
			name:     "超过3个词",
			original: "this is a new method",
			polished: "this is a novel approach",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifier.isSynonymReplacement(tt.original, tt.polished)
			if got != tt.want {
				t.Errorf("isSynonymReplacement() = %v, 期望 %v", got, tt.want)
			}
		})
	}
}

func TestChangeClassifier_isGrammarFix(t *testing.T) {
	classifier := NewChangeClassifier()

	tests := []struct {
		name     string
		original string
		polished string
		want     bool
	}{
		{
			name:     "冠词修正",
			original: "a apple",
			polished: "an apple",
			want:     true,
		},
		{
			name:     "be动词修正",
			original: "he are",
			polished: "he is",
			want:     true,
		},
		{
			name:     "介词修正",
			original: "in Monday",
			polished: "on Monday",
			want:     true,
		},
		{
			name:     "无语法关键词",
			original: "method",
			polished: "approach",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifier.isGrammarFix(tt.original, tt.polished)
			if got != tt.want {
				t.Errorf("isGrammarFix() = %v, 期望 %v", got, tt.want)
			}
		})
	}
}

func TestChangeClassifier_isStructuralChange(t *testing.T) {
	classifier := NewChangeClassifier()

	tests := []struct {
		name     string
		original string
		polished string
		want     bool
	}{
		{
			name:     "词数增加超过30%",
			original: "Hello World",
			polished: "Hello World This is a test",
			want:     true,
		},
		{
			name:     "词数减少超过30%",
			original: "Hello World This is a test",
			polished: "Hello World",
			want:     true,
		},
		{
			name:     "词数变化小于30%",
			original: "Hello World Test",
			polished: "Hello World Check",
			want:     false,
		},
		{
			name:     "相同词数",
			original: "Hello World",
			polished: "Hi Earth",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifier.isStructuralChange(tt.original, tt.polished)
			if got != tt.want {
				t.Errorf("isStructuralChange() = %v, 期望 %v", got, tt.want)
			}
		})
	}
}

func TestChangeClassifier_SuggestHighlightColor(t *testing.T) {
	classifier := NewChangeClassifier()

	tests := []struct {
		name       string
		changeType model.ChangeType
		wantColor  string
	}{
		{
			name:       "词汇优化",
			changeType: model.ChangeTypeVocabulary,
			wantColor:  "yellow",
		},
		{
			name:       "语法修正",
			changeType: model.ChangeTypeGrammar,
			wantColor:  "lightblue",
		},
		{
			name:       "结构调整",
			changeType: model.ChangeTypeStructure,
			wantColor:  "lightgreen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotColor := classifier.SuggestHighlightColor(tt.changeType)
			if gotColor != tt.wantColor {
				t.Errorf("SuggestHighlightColor() = %v, 期望 %v", gotColor, tt.wantColor)
			}
		})
	}
}

func TestChangeClassifier_RealExamples(t *testing.T) {
	classifier := NewChangeClassifier()

	// 真实的润色案例
	examples := []struct {
		name         string
		original     string
		polished     string
		expectedType model.ChangeType
	}{
		{
			name:         "学术词汇提升",
			original:     "use",
			polished:     "utilize",
			expectedType: model.ChangeTypeVocabulary,
		},
		{
			name:         "动词时态修正",
			original:     "he have",
			polished:     "he has",
			expectedType: model.ChangeTypeGrammar,
		},
		{
			name:         "句子重组",
			original:     "the quick brown fox",
			polished:     "a swift, brown fox that moves quickly through the forest",
			expectedType: model.ChangeTypeStructure,
		},
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			gotType := classifier.Classify(ex.original, ex.polished)
			if gotType != ex.expectedType {
				t.Logf("原文: '%s'", ex.original)
				t.Logf("润色: '%s'", ex.polished)
				t.Logf("期望类型: %s, 得到: %s", ex.expectedType, gotType)
			}
		})
	}
}
