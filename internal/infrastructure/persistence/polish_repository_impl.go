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

// polishRepositoryImpl 润色记录仓储实现
type polishRepositoryImpl struct {
	db *gorm.DB
}

// NewPolishRepository 创建润色记录仓储实现
func NewPolishRepository(db *gorm.DB) repository.PolishRepository {
	return &polishRepositoryImpl{db: db}
}

// Create 创建记录
func (r *polishRepositoryImpl) Create(ctx context.Context, record *entity.PolishRecord) error {
	po := &PolishRecordPO{}
	po.FromEntity(record)

	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		logger.Error("failed to create polish record", zap.Error(err))
		return fmt.Errorf("failed to create polish record: %w", err)
	}

	// 回写ID
	record.ID = po.ID
	record.CreatedAt = po.CreatedAt
	record.UpdatedAt = po.UpdatedAt

	return nil
}

// GetByID 根据ID获取记录
func (r *polishRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.PolishRecord, error) {
	var po PolishRecordPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("polish record not found: id=%d", id)
		}
		logger.Error("failed to get polish record by id", zap.Int64("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get polish record: %w", err)
	}

	return po.ToEntity(), nil
}

// GetByTraceID 根据TraceID获取记录
func (r *polishRepositoryImpl) GetByTraceID(ctx context.Context, traceID string) (*entity.PolishRecord, error) {
	var po PolishRecordPO
	err := r.db.WithContext(ctx).Where("trace_id = ?", traceID).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("polish record not found: trace_id=%s", traceID)
		}
		logger.Error("failed to get polish record by trace_id", zap.String("trace_id", traceID), zap.Error(err))
		return nil, fmt.Errorf("failed to get polish record: %w", err)
	}

	return po.ToEntity(), nil
}

// Update 更新记录
func (r *polishRepositoryImpl) Update(ctx context.Context, record *entity.PolishRecord) error {
	po := &PolishRecordPO{}
	po.FromEntity(record)

	result := r.db.WithContext(ctx).Model(&PolishRecordPO{}).Where("id = ?", record.ID).Updates(po)
	if result.Error != nil {
		logger.Error("failed to update polish record", zap.Int64("id", record.ID), zap.Error(result.Error))
		return fmt.Errorf("failed to update polish record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish record not found: id=%d", record.ID)
	}

	return nil
}

// Delete 删除记录（软删除）
func (r *polishRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&PolishRecordPO{}, id)
	if result.Error != nil {
		logger.Error("failed to delete polish record", zap.Int64("id", id), zap.Error(result.Error))
		return fmt.Errorf("failed to delete polish record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("polish record not found: id=%d", id)
	}

	return nil
}

// List 查询记录列表
func (r *polishRepositoryImpl) List(ctx context.Context, opts repository.QueryOptions) ([]*entity.PolishRecord, error) {
	var pos []*PolishRecordPO

	query := r.buildQuery(ctx, opts)

	// 排序
	if opts.OrderBy != "" {
		order := opts.OrderBy
		if opts.OrderDesc {
			order += " DESC"
		}
		query = query.Order(order)
	}

	// 分页
	if opts.Limit > 0 {
		query = query.Offset(opts.Offset).Limit(opts.Limit)
	}

	// 执行查询
	if err := query.Find(&pos).Error; err != nil {
		logger.Error("failed to list polish records", zap.Error(err))
		return nil, fmt.Errorf("failed to list polish records: %w", err)
	}

	return ToEntityList(pos), nil
}

// Count 统计记录数量
func (r *polishRepositoryImpl) Count(ctx context.Context, opts repository.QueryOptions) (int64, error) {
	var count int64

	query := r.buildQuery(ctx, opts)

	if err := query.Model(&PolishRecordPO{}).Count(&count).Error; err != nil {
		logger.Error("failed to count polish records", zap.Error(err))
		return 0, fmt.Errorf("failed to count polish records: %w", err)
	}

	return count, nil
}

// BatchCreate 批量创建记录
func (r *polishRepositoryImpl) BatchCreate(ctx context.Context, records []*entity.PolishRecord) error {
	pos := make([]*PolishRecordPO, len(records))
	for i, record := range records {
		po := &PolishRecordPO{}
		po.FromEntity(record)
		pos[i] = po
	}

	if err := r.db.WithContext(ctx).Create(&pos).Error; err != nil {
		logger.Error("failed to batch create polish records", zap.Error(err))
		return fmt.Errorf("failed to batch create polish records: %w", err)
	}

	// 回写ID
	for i, po := range pos {
		records[i].ID = po.ID
		records[i].CreatedAt = po.CreatedAt
		records[i].UpdatedAt = po.UpdatedAt
	}

	return nil
}

// buildQuery 构建查询条件
func (r *polishRepositoryImpl) buildQuery(ctx context.Context, opts repository.QueryOptions) *gorm.DB {
	query := r.db.WithContext(ctx)

	// 字段选择优化
	if opts.ExcludeText {
		// 排除大文本字段，提高查询性能
		query = query.Select("id, trace_id, user_id, style, language, original_length, polished_length, provider, model, process_time_ms, status, created_at, updated_at")
	} else if len(opts.SelectFields) > 0 {
		query = query.Select(opts.SelectFields)
	}

	// 过滤条件
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}

	if opts.Provider != nil {
		query = query.Where("provider = ?", *opts.Provider)
	}

	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	if opts.Language != nil {
		query = query.Where("language = ?", *opts.Language)
	}

	if opts.Style != nil {
		query = query.Where("style = ?", *opts.Style)
	}

	// 时间范围过滤
	if opts.StartTime != nil {
		query = query.Where("created_at >= ?", *opts.StartTime)
	}

	if opts.EndTime != nil {
		query = query.Where("created_at <= ?", *opts.EndTime)
	}

	return query
}
