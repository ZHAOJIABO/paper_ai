package repository

import (
	"context"
	"paper_ai/internal/domain/entity"
)

// RefreshTokenRepository 刷新令牌仓储接口
type RefreshTokenRepository interface {
	// Create 创建刷新令牌
	Create(ctx context.Context, token *entity.RefreshToken) error

	// GetByToken 根据token字符串获取令牌
	GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error)

	// Revoke 撤销令牌
	Revoke(ctx context.Context, token string) error

	// RevokeAllByUserID 撤销用户的所有令牌（用户登出所有设备）
	RevokeAllByUserID(ctx context.Context, userID int64) error

	// DeleteExpired 删除过期的令牌（定期清理）
	DeleteExpired(ctx context.Context) error

	// GetValidTokensByUserID 获取用户的有效令牌列表
	GetValidTokensByUserID(ctx context.Context, userID int64) ([]*entity.RefreshToken, error)
}
