package persistence

import (
	"context"
	"fmt"

	"paper_ai/internal/domain/repository"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
)

// GetStatistics 获取统计信息
func (r *polishRepositoryImpl) GetStatistics(ctx context.Context, opts repository.StatisticsOptions) (*repository.Statistics, error) {
	stats := &repository.Statistics{
		ProviderStats: make(map[string]*repository.Stat),
		LanguageStats: make(map[string]*repository.Stat),
		StyleStats:    make(map[string]*repository.Stat),
	}

	query := r.db.WithContext(ctx).Model(&PolishRecordPO{})

	// 用户ID过滤
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}

	// 时间范围过滤
	if opts.TimeRange != nil {
		query = query.Where("created_at >= ? AND created_at <= ?", opts.TimeRange.Start, opts.TimeRange.End)
	}

	// 总记录数
	if err := query.Count(&stats.TotalCount).Error; err != nil {
		logger.Error("failed to count total records", zap.Error(err))
		return nil, fmt.Errorf("failed to count total records: %w", err)
	}

	// 成功记录数
	if err := query.Where("status = ?", "success").Count(&stats.SuccessCount).Error; err != nil {
		logger.Error("failed to count success records", zap.Error(err))
		return nil, fmt.Errorf("failed to count success records: %w", err)
	}

	// 失败记录数
	stats.FailedCount = stats.TotalCount - stats.SuccessCount

	// 成功率
	if stats.TotalCount > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.TotalCount) * 100
	}

	// 平均处理时间（仅成功的记录）
	var avgTime float64
	if err := query.Where("status = ?", "success").Select("AVG(process_time_ms)").Scan(&avgTime).Error; err != nil {
		logger.Error("failed to calculate avg process time", zap.Error(err))
	} else {
		stats.AvgProcessTime = avgTime
	}

	// 按提供商统计
	if err := r.getProviderStats(ctx, opts, stats); err != nil {
		logger.Warn("failed to get provider stats", zap.Error(err))
	}

	// 按语言统计
	if err := r.getLanguageStats(ctx, opts, stats); err != nil {
		logger.Warn("failed to get language stats", zap.Error(err))
	}

	// 按风格统计
	if err := r.getStyleStats(ctx, opts, stats); err != nil {
		logger.Warn("failed to get style stats", zap.Error(err))
	}

	return stats, nil
}

// getProviderStats 按提供商统计
func (r *polishRepositoryImpl) getProviderStats(ctx context.Context, opts repository.StatisticsOptions, stats *repository.Statistics) error {
	type ProviderStat struct {
		Provider       string
		Count          int64
		SuccessCount   int64
		AvgProcessTime float64
	}

	query := r.db.WithContext(ctx).Model(&PolishRecordPO{})
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.TimeRange != nil {
		query = query.Where("created_at >= ? AND created_at <= ?", opts.TimeRange.Start, opts.TimeRange.End)
	}

	var providerStats []ProviderStat
	if err := query.Select(
		"provider",
		"COUNT(*) as count",
		"SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_count",
		"AVG(CASE WHEN status = 'success' THEN process_time_ms ELSE 0 END) as avg_process_time",
	).Group("provider").Scan(&providerStats).Error; err != nil {
		return err
	}

	for _, ps := range providerStats {
		stat := &repository.Stat{
			Count:          ps.Count,
			SuccessCount:   ps.SuccessCount,
			FailedCount:    ps.Count - ps.SuccessCount,
			AvgProcessTime: ps.AvgProcessTime,
		}
		if ps.Count > 0 {
			stat.SuccessRate = float64(ps.SuccessCount) / float64(ps.Count) * 100
		}
		stats.ProviderStats[ps.Provider] = stat
	}

	return nil
}

// getLanguageStats 按语言统计
func (r *polishRepositoryImpl) getLanguageStats(ctx context.Context, opts repository.StatisticsOptions, stats *repository.Statistics) error {
	type LanguageStat struct {
		Language       string
		Count          int64
		SuccessCount   int64
		AvgProcessTime float64
	}

	query := r.db.WithContext(ctx).Model(&PolishRecordPO{})
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.TimeRange != nil {
		query = query.Where("created_at >= ? AND created_at <= ?", opts.TimeRange.Start, opts.TimeRange.End)
	}

	var languageStats []LanguageStat
	if err := query.Select(
		"language",
		"COUNT(*) as count",
		"SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_count",
		"AVG(CASE WHEN status = 'success' THEN process_time_ms ELSE 0 END) as avg_process_time",
	).Group("language").Scan(&languageStats).Error; err != nil {
		return err
	}

	for _, ls := range languageStats {
		stat := &repository.Stat{
			Count:          ls.Count,
			SuccessCount:   ls.SuccessCount,
			FailedCount:    ls.Count - ls.SuccessCount,
			AvgProcessTime: ls.AvgProcessTime,
		}
		if ls.Count > 0 {
			stat.SuccessRate = float64(ls.SuccessCount) / float64(ls.Count) * 100
		}
		stats.LanguageStats[ls.Language] = stat
	}

	return nil
}

// getStyleStats 按风格统计
func (r *polishRepositoryImpl) getStyleStats(ctx context.Context, opts repository.StatisticsOptions, stats *repository.Statistics) error {
	type StyleStat struct {
		Style          string
		Count          int64
		SuccessCount   int64
		AvgProcessTime float64
	}

	query := r.db.WithContext(ctx).Model(&PolishRecordPO{})
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.TimeRange != nil {
		query = query.Where("created_at >= ? AND created_at <= ?", opts.TimeRange.Start, opts.TimeRange.End)
	}

	var styleStats []StyleStat
	if err := query.Select(
		"style",
		"COUNT(*) as count",
		"SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_count",
		"AVG(CASE WHEN status = 'success' THEN process_time_ms ELSE 0 END) as avg_process_time",
	).Group("style").Scan(&styleStats).Error; err != nil {
		return err
	}

	for _, ss := range styleStats {
		stat := &repository.Stat{
			Count:          ss.Count,
			SuccessCount:   ss.SuccessCount,
			FailedCount:    ss.Count - ss.SuccessCount,
			AvgProcessTime: ss.AvgProcessTime,
		}
		if ss.Count > 0 {
			stat.SuccessRate = float64(ss.SuccessCount) / float64(ss.Count) * 100
		}
		stats.StyleStats[ss.Style] = stat
	}

	return nil
}
