package persistence

import (
	"context"
	"fmt"

	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/repository"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// polishVersionRepositoryImpl 润色版本仓储实现
type polishVersionRepositoryImpl struct {
	db *gorm.DB
}

// NewPolishVersionRepository 创建润色版本仓储实现
func NewPolishVersionRepository(db *gorm.DB) repository.PolishVersionRepository {
	return &polishVersionRepositoryImpl{db: db}
}

// Create 创建版本记录
func (r *polishVersionRepositoryImpl) Create(ctx context.Context, version *entity.PolishVersion) error {
	po := &PolishVersionPO{}
	po.FromEntity(version)

	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		logger.Error("failed to create polish version", zap.Error(err))
		return fmt.Errorf("failed to create polish version: %w", err)
	}

	// 回写ID和时间戳
	version.ID = po.ID
	version.CreatedAt = po.CreatedAt

	return nil
}

// CreateBatch 批量创建版本记录
func (r *polishVersionRepositoryImpl) CreateBatch(ctx context.Context, versions []*entity.PolishVersion) error {
	if len(versions) == 0 {
		return nil
	}

	// 转换为PO
	pos := make([]*PolishVersionPO, len(versions))
	for i, v := range versions {
		po := &PolishVersionPO{}
		po.FromEntity(v)
		pos[i] = po
	}

	// 批量插入
	if err := r.db.WithContext(ctx).Create(&pos).Error; err != nil {
		logger.Error("failed to batch create polish versions", zap.Error(err))
		return fmt.Errorf("failed to batch create polish versions: %w", err)
	}

	// 回写ID和时间戳
	for i, po := range pos {
		versions[i].ID = po.ID
		versions[i].CreatedAt = po.CreatedAt
	}

	return nil
}

// GetByID 根据ID获取版本
func (r *polishVersionRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.PolishVersion, error) {
	var po PolishVersionPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("polish version not found: id=%d", id)
		}
		logger.Error("failed to get polish version by id", zap.Int64("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get polish version: %w", err)
	}

	return po.ToEntity(), nil
}

// GetByRecordID 获取某条主记录的所有版本
func (r *polishVersionRepositoryImpl) GetByRecordID(ctx context.Context, recordID int64) ([]*entity.PolishVersion, error) {
	var pos []*PolishVersionPO
	err := r.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Order("created_at ASC").
		Find(&pos).Error

	if err != nil {
		logger.Error("failed to get polish versions by record_id", zap.Int64("record_id", recordID), zap.Error(err))
		return nil, fmt.Errorf("failed to get polish versions: %w", err)
	}

	// 转换为实体
	versions := make([]*entity.PolishVersion, len(pos))
	for i, po := range pos {
		versions[i] = po.ToEntity()
	}

	return versions, nil
}

// GetByRecordIDAndType 获取某条主记录的特定版本类型
func (r *polishVersionRepositoryImpl) GetByRecordIDAndType(ctx context.Context, recordID int64, versionType string) (*entity.PolishVersion, error) {
	var po PolishVersionPO
	err := r.db.WithContext(ctx).
		Where("record_id = ? AND version_type = ?", recordID, versionType).
		First(&po).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("polish version not found: record_id=%d, version_type=%s", recordID, versionType)
		}
		logger.Error("failed to get polish version", zap.Int64("record_id", recordID), zap.String("version_type", versionType), zap.Error(err))
		return nil, fmt.Errorf("failed to get polish version: %w", err)
	}

	return po.ToEntity(), nil
}

// Update 更新版本记录
func (r *polishVersionRepositoryImpl) Update(ctx context.Context, version *entity.PolishVersion) error {
	po := &PolishVersionPO{}
	po.FromEntity(version)

	result := r.db.WithContext(ctx).Model(&PolishVersionPO{}).Where("id = ?", version.ID).Updates(po)
	if result.Error != nil {
		logger.Error("failed to update polish version", zap.Int64("id", version.ID), zap.Error(result.Error))
		return fmt.Errorf("failed to update polish version: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish version not found: id=%d", version.ID)
	}

	return nil
}

// Delete 删除版本记录
func (r *polishVersionRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&PolishVersionPO{}, id)
	if result.Error != nil {
		logger.Error("failed to delete polish version", zap.Int64("id", id), zap.Error(result.Error))
		return fmt.Errorf("failed to delete polish version: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish version not found: id=%d", id)
	}

	return nil
}

// DeleteByRecordID 删除某条主记录的所有版本
func (r *polishVersionRepositoryImpl) DeleteByRecordID(ctx context.Context, recordID int64) error {
	result := r.db.WithContext(ctx).Where("record_id = ?", recordID).Delete(&PolishVersionPO{})
	if result.Error != nil {
		logger.Error("failed to delete polish versions by record_id", zap.Int64("record_id", recordID), zap.Error(result.Error))
		return fmt.Errorf("failed to delete polish versions: %w", result.Error)
	}

	return nil
}

// Count 统计版本数量
func (r *polishVersionRepositoryImpl) Count(ctx context.Context, filter repository.VersionFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&PolishVersionPO{})

	// 应用过滤条件
	if filter.RecordID != 0 {
		query = query.Where("record_id = ?", filter.RecordID)
	}
	if filter.VersionType != "" {
		query = query.Where("version_type = ?", filter.VersionType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		logger.Error("failed to count polish versions", zap.Error(err))
		return 0, fmt.Errorf("failed to count polish versions: %w", err)
	}

	return count, nil
}

// GetStatsByVersionType 按版本类型统计
func (r *polishVersionRepositoryImpl) GetStatsByVersionType(ctx context.Context) (map[string]*repository.VersionTypeStats, error) {
	type statsResult struct {
		VersionType      string
		TotalCount       int64
		SuccessCount     int64
		FailedCount      int64
		AvgProcessTimeMs float64
	}

	var results []statsResult
	err := r.db.WithContext(ctx).Model(&PolishVersionPO{}).
		Select(`
			version_type,
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_count,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed_count,
			AVG(process_time_ms) as avg_process_time_ms
		`).
		Group("version_type").
		Find(&results).Error

	if err != nil {
		logger.Error("failed to get stats by version type", zap.Error(err))
		return nil, fmt.Errorf("failed to get stats by version type: %w", err)
	}

	// 转换为map
	statsMap := make(map[string]*repository.VersionTypeStats)
	for _, r := range results {
		statsMap[r.VersionType] = &repository.VersionTypeStats{
			VersionType:      r.VersionType,
			TotalCount:       r.TotalCount,
			SuccessCount:     r.SuccessCount,
			FailedCount:      r.FailedCount,
			AvgProcessTimeMs: r.AvgProcessTimeMs,
		}
	}

	return statsMap, nil
}
