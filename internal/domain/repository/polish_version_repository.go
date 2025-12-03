package repository

import (
	"context"
	"paper_ai/internal/domain/entity"
)

// PolishVersionRepository 润色版本仓储接口
type PolishVersionRepository interface {
	// Create 创建版本记录
	Create(ctx context.Context, version *entity.PolishVersion) error

	// CreateBatch 批量创建版本记录
	CreateBatch(ctx context.Context, versions []*entity.PolishVersion) error

	// GetByID 根据ID获取版本
	GetByID(ctx context.Context, id int64) (*entity.PolishVersion, error)

	// GetByRecordID 获取某条主记录的所有版本
	GetByRecordID(ctx context.Context, recordID int64) ([]*entity.PolishVersion, error)

	// GetByRecordIDAndType 获取某条主记录的特定版本类型
	GetByRecordIDAndType(ctx context.Context, recordID int64, versionType string) (*entity.PolishVersion, error)

	// Update 更新版本记录
	Update(ctx context.Context, version *entity.PolishVersion) error

	// Delete 删除版本记录
	Delete(ctx context.Context, id int64) error

	// DeleteByRecordID 删除某条主记录的所有版本
	DeleteByRecordID(ctx context.Context, recordID int64) error

	// Count 统计版本数量
	Count(ctx context.Context, filter VersionFilter) (int64, error)

	// GetStatsByVersionType 按版本类型统计
	GetStatsByVersionType(ctx context.Context) (map[string]*VersionTypeStats, error)
}

// VersionFilter 版本查询过滤器
type VersionFilter struct {
	RecordID    int64
	VersionType string
	Status      string
}

// VersionTypeStats 版本类型统计
type VersionTypeStats struct {
	VersionType      string
	TotalCount       int64
	SuccessCount     int64
	FailedCount      int64
	AvgProcessTimeMs float64
}
