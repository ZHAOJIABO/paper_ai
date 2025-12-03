package service

import (
	"context"
	"fmt"

	"paper_ai/internal/domain/repository"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
)

// FeatureService 功能开关服务
type FeatureService struct {
	userRepo repository.UserRepository
	config   *FeatureConfig
}

// FeatureConfig 功能配置
type FeatureConfig struct {
	// 多版本润色功能全局开关
	MultiVersionEnabled bool
	// 默认模式
	DefaultMode string
	// 最大并发数
	MaxConcurrent int
}

// NewFeatureService 创建功能开关服务
func NewFeatureService(userRepo repository.UserRepository, config *FeatureConfig) *FeatureService {
	return &FeatureService{
		userRepo: userRepo,
		config:   config,
	}
}

// CheckMultiVersionPermission 检查用户是否有多版本润色权限
// 返回：hasPermission, reason, error
func (s *FeatureService) CheckMultiVersionPermission(ctx context.Context, userID int64) (bool, string, error) {
	// 1. 检查全局开关
	if !s.config.MultiVersionEnabled {
		logger.Warn("multi-version feature is globally disabled")
		return false, "多版本润色功能暂未开放", nil
	}

	// 2. 查询用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Error("failed to get user info", zap.Int64("user_id", userID), zap.Error(err))
		return false, "", fmt.Errorf("failed to get user info: %w", err)
	}

	// 3. 检查用户权限
	if !user.HasMultiVersionPermission() {
		logger.Warn("user does not have multi-version permission", zap.Int64("user_id", userID))
		return false, "您暂无使用多版本润色的权限，请联系管理员开通", nil
	}

	// 4. 检查配额（如果设置了配额）
	if !user.HasUnlimitedQuota() {
		// TODO: 查询用户已使用的配额
		// 这里需要查询 polish_records 表，统计该用户的多版本润色次数
		// 暂时跳过配额检查，后续实现
		logger.Debug("user has quota limit", zap.Int64("user_id", userID), zap.Int("quota", user.MultiVersionQuota))
	}

	logger.Info("user has multi-version permission", zap.Int64("user_id", userID))
	return true, "", nil
}

// IsMultiVersionEnabled 判断多版本功能是否全局启用
func (s *FeatureService) IsMultiVersionEnabled() bool {
	return s.config.MultiVersionEnabled
}

// GetDefaultMode 获取默认模式
func (s *FeatureService) GetDefaultMode() string {
	return s.config.DefaultMode
}

// GetMaxConcurrent 获取最大并发数
func (s *FeatureService) GetMaxConcurrent() int {
	return s.config.MaxConcurrent
}
