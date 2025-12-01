package comparison

import (
	"strings"
	"testing"

	"paper_ai/internal/domain/model"
)

func TestReasonGenerator_Generate(t *testing.T) {
	generator := NewReasonGenerator()

	tests := []struct {
		name       string
		changeType model.ChangeType
		original   string
		polished   string
		wantContain string // 期望理由中包含的关键词
	}{
		{
			name:        "词汇优化理由",
			changeType:  model.ChangeTypeVocabulary,
			original:    "method",
			polished:    "methodology",
			wantContain: "学术",
		},
		{
			name:        "语法优化理由",
			changeType:  model.ChangeTypeGrammar,
			original:    "a apple",
			polished:    "an apple",
			wantContain: "语法",
		},
		{
			name:        "结构调整理由",
			changeType:  model.ChangeTypeStructure,
			original:    "short",
			polished:    "this is a much longer sentence",
			wantContain: "结构",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reason := generator.Generate(tt.changeType, tt.original, tt.polished)

			if reason == "" {
				t.Error("Generate() 返回空字符串")
			}

			if !strings.Contains(reason, tt.wantContain) {
				t.Errorf("Generate() 理由中应包含 '%s', 实际: %s", tt.wantContain, reason)
			}

			t.Logf("修改类型: %s", tt.changeType)
			t.Logf("生成理由: %s", reason)
		})
	}
}

func TestReasonGenerator_GenerateAlternatives(t *testing.T) {
	generator := NewReasonGenerator()

	tests := []struct {
		name         string
		changeType   model.ChangeType
		original     string
		wantMinCount int // 期望至少有几个替代方案
	}{
		{
			name:         "词汇优化有替代方案",
			changeType:   model.ChangeTypeVocabulary,
			original:     "method",
			wantMinCount: 1,
		},
		{
			name:         "语法修正无替代方案",
			changeType:   model.ChangeTypeGrammar,
			original:     "is",
			wantMinCount: 0,
		},
		{
			name:         "结构调整无替代方案",
			changeType:   model.ChangeTypeStructure,
			original:     "sentence",
			wantMinCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alternatives := generator.GenerateAlternatives(tt.changeType, tt.original)

			if len(alternatives) < tt.wantMinCount {
				t.Errorf("GenerateAlternatives() 返回 %d 个方案，期望至少 %d 个", len(alternatives), tt.wantMinCount)
			}

			for i, alt := range alternatives {
				t.Logf("替代方案 %d: '%s' - %s", i+1, alt.Text, alt.Reason)

				if alt.Text == "" {
					t.Error("替代方案文本不能为空")
				}
				if alt.Reason == "" {
					t.Error("替代方案理由不能为空")
				}
			}
		})
	}
}

func TestReasonGenerator_VocabularyAlternatives(t *testing.T) {
	generator := NewReasonGenerator()

	// 测试常见词汇的替代方案
	commonWords := []string{"method", "problem", "solve", "show", "use", "get", "make", "help"}

	for _, word := range commonWords {
		t.Run("词汇_"+word, func(t *testing.T) {
			alternatives := generator.getVocabularyAlternatives(word)

			if len(alternatives) == 0 {
				t.Errorf("常见词汇 '%s' 应该有替代方案", word)
			}

			for _, alt := range alternatives {
				t.Logf("'%s' 的替代方案: '%s' - %s", word, alt.Text, alt.Reason)
			}
		})
	}
}

func TestReasonGenerator_CalculateConfidence(t *testing.T) {
	generator := NewReasonGenerator()

	tests := []struct {
		name       string
		changeType model.ChangeType
		wantMin    float64
		wantMax    float64
	}{
		{
			name:       "词汇优化置信度高",
			changeType: model.ChangeTypeVocabulary,
			wantMin:    0.85,
			wantMax:    1.0,
		},
		{
			name:       "语法修正置信度中等",
			changeType: model.ChangeTypeGrammar,
			wantMin:    0.80,
			wantMax:    0.95,
		},
		{
			name:       "结构调整置信度较低",
			changeType: model.ChangeTypeStructure,
			wantMin:    0.70,
			wantMax:    0.85,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confidence := generator.CalculateConfidence(tt.changeType, "original", "polished")

			if confidence < tt.wantMin || confidence > tt.wantMax {
				t.Errorf("CalculateConfidence() = %v, 期望在 [%v, %v] 范围内", confidence, tt.wantMin, tt.wantMax)
			}

			if confidence < 0 || confidence > 1 {
				t.Errorf("置信度应该在 [0, 1] 范围内，得到 %v", confidence)
			}

			t.Logf("修改类型: %s, 置信度: %.2f", tt.changeType, confidence)
		})
	}
}

func TestReasonGenerator_GetImpact(t *testing.T) {
	generator := NewReasonGenerator()

	tests := []struct {
		name       string
		changeType model.ChangeType
		wantImpact string
	}{
		{
			name:       "词汇优化影响学术语气",
			changeType: model.ChangeTypeVocabulary,
			wantImpact: "academic_tone",
		},
		{
			name:       "语法修正影响语法正确性",
			changeType: model.ChangeTypeGrammar,
			wantImpact: "grammar_correctness",
		},
		{
			name:       "结构调整影响可读性",
			changeType: model.ChangeTypeStructure,
			wantImpact: "readability",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impact := generator.GetImpact(tt.changeType)

			if impact != tt.wantImpact {
				t.Errorf("GetImpact() = %v, 期望 %v", impact, tt.wantImpact)
			}
		})
	}
}

func TestReasonGenerator_CompleteFunctionality(t *testing.T) {
	generator := NewReasonGenerator()

	// 测试完整的生成流程
	testCases := []struct {
		changeType model.ChangeType
		original   string
		polished   string
	}{
		{model.ChangeTypeVocabulary, "method", "methodology"},
		{model.ChangeTypeGrammar, "he have", "he has"},
		{model.ChangeTypeStructure, "short", "this is a longer sentence"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.changeType), func(t *testing.T) {
			// 生成理由
			reason := generator.Generate(tc.changeType, tc.original, tc.polished)
			if reason == "" {
				t.Error("理由生成失败")
			}

			// 生成替代方案
			alternatives := generator.GenerateAlternatives(tc.changeType, tc.original)
			t.Logf("替代方案数量: %d", len(alternatives))

			// 计算置信度
			confidence := generator.CalculateConfidence(tc.changeType, tc.original, tc.polished)
			if confidence < 0 || confidence > 1 {
				t.Errorf("置信度超出范围: %v", confidence)
			}

			// 获取影响维度
			impact := generator.GetImpact(tc.changeType)
			if impact == "" {
				t.Error("影响维度为空")
			}

			t.Logf("修改: '%s' -> '%s'", tc.original, tc.polished)
			t.Logf("理由: %s", reason)
			t.Logf("置信度: %.2f", confidence)
			t.Logf("影响: %s", impact)
		})
	}
}
