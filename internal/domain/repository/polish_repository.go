package repository

import (
	"context"
	"paper_ai/internal/domain/entity"
)

// PolishRepository 润色记录仓储接口
// 定义所有数据访问操作的契约，与具体实现无关（依赖倒置原则）
type PolishRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, record *entity.PolishRecord) error
	GetByID(ctx context.Context, id int64) (*entity.PolishRecord, error)
	GetByTraceID(ctx context.Context, traceID string) (*entity.PolishRecord, error)
	Update(ctx context.Context, record *entity.PolishRecord) error
	Delete(ctx context.Context, id int64) error

	// 查询操作（使用Options模式，高度可扩展）
	List(ctx context.Context, opts QueryOptions) ([]*entity.PolishRecord, error)
	Count(ctx context.Context, opts QueryOptions) (int64, error)

	// 统计操作
	GetStatistics(ctx context.Context, opts StatisticsOptions) (*Statistics, error)

	// 批量操作（可选，用于未来扩展）
	BatchCreate(ctx context.Context, records []*entity.PolishRecord) error
}
