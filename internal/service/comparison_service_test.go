package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/repository"
	"paper_ai/pkg/logger"
)

func init() {
	// 初始化logger for 测试
	_ = logger.Init()
}

// MockPolishRepository 模拟润色记录仓储
type MockPolishRepository struct {
	records map[string]*entity.PolishRecord
}

func NewMockPolishRepository() *MockPolishRepository {
	return &MockPolishRepository{
		records: make(map[string]*entity.PolishRecord),
	}
}

func (m *MockPolishRepository) GetByTraceID(ctx context.Context, traceID string) (*entity.PolishRecord, error) {
	if record, ok := m.records[traceID]; ok {
		return record, nil
	}
	return nil, fmt.Errorf("记录不存在: trace_id=%s", traceID)
}

func (m *MockPolishRepository) Update(ctx context.Context, record *entity.PolishRecord) error {
	m.records[record.TraceID] = record
	return nil
}

// 添加模拟记录
func (m *MockPolishRepository) AddMockRecord(record *entity.PolishRecord) {
	m.records[record.TraceID] = record
}

// 实现其他必需的接口方法（空实现）
func (m *MockPolishRepository) Create(ctx context.Context, record *entity.PolishRecord) error {
	return nil
}
func (m *MockPolishRepository) GetByID(ctx context.Context, id int64) (*entity.PolishRecord, error) {
	return nil, nil
}
func (m *MockPolishRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
func (m *MockPolishRepository) List(ctx context.Context, opts repository.QueryOptions) ([]*entity.PolishRecord, error) {
	return nil, nil
}
func (m *MockPolishRepository) Count(ctx context.Context, opts repository.QueryOptions) (int64, error) {
	return 0, nil
}
func (m *MockPolishRepository) BatchCreate(ctx context.Context, records []*entity.PolishRecord) error {
	return nil
}
func (m *MockPolishRepository) GetStatistics(ctx context.Context, opts repository.StatisticsOptions) (*repository.Statistics, error) {
	return nil, nil
}

func TestComparisonService_GenerateComparison(t *testing.T) {
	mockRepo := NewMockPolishRepository()
	service := NewComparisonService(mockRepo)

	// 准备测试数据
	testRecord := &entity.PolishRecord{
		TraceID:         "1732701603123",
		UserID:          12345,
		OriginalContent: "In this paper, we propose a new method to solve the problem.",
		PolishedContent: "In this paper, we propose a novel methodology to address the issue.",
		Status:          "success",
		CreatedAt:       time.Now(),
	}
	mockRepo.AddMockRecord(testRecord)

	ctx := context.Background()

	t.Run("生成新的对比数据", func(t *testing.T) {
		result, err := service.GenerateComparison(ctx, "1732701603123")
		if err != nil {
			t.Fatalf("GenerateComparison() 失败: %v", err)
		}

		// 验证基本信息
		if result.TraceID != "1732701603123" {
			t.Errorf("TraceID = %s, 期望 1732701603123", result.TraceID)
		}

		if result.OriginalContent == "" || result.PolishedContent == "" {
			t.Error("原文或润色文本为空")
		}

		// 验证标注列表
		if len(result.Annotations) == 0 {
			t.Error("应该检测到至少一处修改")
		}

		t.Logf("检测到 %d 处修改:", len(result.Annotations))
		for i, ann := range result.Annotations {
			t.Logf("  修改 %d: '%s' -> '%s' (类型: %s)", i+1, ann.OriginalText, ann.PolishedText, ann.Type)
			t.Logf("    位置: [%d:%d]", ann.PolishedPosition.Start, ann.PolishedPosition.End)
			t.Logf("    理由: %s", ann.Reason)
			t.Logf("    置信度: %.2f", ann.Confidence)

			// 验证必填字段
			if ann.ID == "" {
				t.Error("修改ID不能为空")
			}
			if ann.Reason == "" {
				t.Error("修改理由不能为空")
			}
			if ann.Confidence < 0 || ann.Confidence > 1 {
				t.Errorf("置信度超出范围: %v", ann.Confidence)
			}
		}

		// 验证元数据
		if result.Metadata.TotalChanges != len(result.Annotations) {
			t.Errorf("元数据总修改数(%d)与标注列表长度(%d)不一致", result.Metadata.TotalChanges, len(result.Annotations))
		}

		t.Logf("元数据:")
		t.Logf("  原文词数: %d", result.Metadata.OriginalWordCount)
		t.Logf("  润色词数: %d", result.Metadata.PolishedWordCount)
		t.Logf("  总修改数: %d", result.Metadata.TotalChanges)
		t.Logf("  学术性提升: %.2f%%", result.Metadata.AcademicScoreImprovement)

		// 验证统计信息
		totalChanges := result.Statistics.VocabularyChanges + result.Statistics.GrammarChanges + result.Statistics.StructureChanges
		if totalChanges != len(result.Annotations) {
			t.Errorf("统计信息总数(%d)与标注列表长度(%d)不一致", totalChanges, len(result.Annotations))
		}

		t.Logf("统计信息:")
		t.Logf("  词汇优化: %d", result.Statistics.VocabularyChanges)
		t.Logf("  语法修正: %d", result.Statistics.GrammarChanges)
		t.Logf("  结构调整: %d", result.Statistics.StructureChanges)
	})

	t.Run("记录不存在", func(t *testing.T) {
		_, err := service.GenerateComparison(ctx, "9999999999999")
		if err == nil {
			t.Error("应该返回错误")
		} else {
			t.Logf("正确返回错误: %v", err)
		}
	})
}

func TestComparisonService_GetComparison(t *testing.T) {
	mockRepo := NewMockPolishRepository()
	service := NewComparisonService(mockRepo)

	testRecord := &entity.PolishRecord{
		TraceID:         "1732701603456",
		UserID:          12345,
		OriginalContent: "This is a test method.",
		PolishedContent: "This is a test methodology.",
		Status:          "success",
		CreatedAt:       time.Now(),
	}
	mockRepo.AddMockRecord(testRecord)

	ctx := context.Background()

	t.Run("获取对比数据", func(t *testing.T) {
		result, err := service.GetComparison(ctx, "1732701603456", 12345)
		if err != nil {
			t.Fatalf("GetComparison() 失败: %v", err)
		}

		if result == nil {
			t.Fatal("返回结果为空")
		}

		t.Logf("成功获取对比数据，检测到 %d 处修改", len(result.Annotations))
	})

	t.Run("权限验证 - 不同用户", func(t *testing.T) {
		_, err := service.GetComparison(ctx, "1732701603456", 99999)
		if err == nil {
			t.Error("应该拒绝不同用户的访问")
		}
	})
}

func TestComparisonService_ComplexText(t *testing.T) {
	mockRepo := NewMockPolishRepository()
	service := NewComparisonService(mockRepo)

	// 更复杂的文本
	complexRecord := &entity.PolishRecord{
		TraceID: "1732701603789",
		UserID:  12345,
		OriginalContent: `In recent years, the use of machine learning has grown rapidly.
Many researchers try to solve difficult problems using these new methods.
However, there are still many challenges to overcome.`,
		PolishedContent: `In recent years, the utilization of machine learning has expanded significantly.
Numerous researchers endeavor to address complex challenges through these innovative methodologies.
Nevertheless, substantial obstacles remain to be surmounted.`,
		Status:    "success",
		CreatedAt: time.Now(),
	}
	mockRepo.AddMockRecord(complexRecord)

	ctx := context.Background()

	result, err := service.GenerateComparison(ctx, "1732701603789")
	if err != nil {
		t.Fatalf("GenerateComparison() 失败: %v", err)
	}

	t.Logf("复杂文本对比分析:")
	t.Logf("原文长度: %d 词", result.Metadata.OriginalWordCount)
	t.Logf("润色长度: %d 词", result.Metadata.PolishedWordCount)
	t.Logf("检测到 %d 处修改", len(result.Annotations))

	if len(result.Annotations) == 0 {
		t.Error("应该检测到多处修改")
	}

	// 按类型分组统计
	vocabCount := 0
	grammarCount := 0
	structureCount := 0

	for _, ann := range result.Annotations {
		switch ann.Type {
		case "vocabulary":
			vocabCount++
		case "grammar":
			grammarCount++
		case "structure":
			structureCount++
		}
	}

	t.Logf("修改类型分布:")
	t.Logf("  词汇优化: %d", vocabCount)
	t.Logf("  语法修正: %d", grammarCount)
	t.Logf("  结构调整: %d", structureCount)
}

func TestComparisonService_SaveAndRetrieve(t *testing.T) {
	mockRepo := NewMockPolishRepository()
	service := NewComparisonService(mockRepo)

	testRecord := &entity.PolishRecord{
		TraceID:         "1732701603999",
		UserID:          12345,
		OriginalContent: "simple text",
		PolishedContent: "complex text",
		Status:          "success",
		CreatedAt:       time.Now(),
	}
	mockRepo.AddMockRecord(testRecord)

	ctx := context.Background()

	// 第一次生成
	result1, err := service.GenerateComparison(ctx, "1732701603999")
	if err != nil {
		t.Fatalf("第一次生成失败: %v", err)
	}

	// 验证数据已保存
	savedRecord, _ := mockRepo.GetByTraceID(ctx, "1732701603999")
	if savedRecord.ComparisonData == "" {
		t.Error("对比数据未保存")
	}

	if savedRecord.ChangesCount != len(result1.Annotations) {
		t.Errorf("保存的修改数量不正确: 期望 %d, 得到 %d", len(result1.Annotations), savedRecord.ChangesCount)
	}

	t.Logf("对比数据已保存，修改数量: %d", savedRecord.ChangesCount)
}
