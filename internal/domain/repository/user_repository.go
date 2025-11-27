package repository

import (
	"context"
	"paper_ai/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id int64) (*entity.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *entity.User) error

	// UpdateLoginInfo 更新登录信息
	UpdateLoginInfo(ctx context.Context, userID int64, ip string) error

	// IncrementFailedLoginCount 增加失败登录次数
	IncrementFailedLoginCount(ctx context.Context, userID int64) error

	// ResetFailedLoginCount 重置失败登录次数
	ResetFailedLoginCount(ctx context.Context, userID int64) error

	// ExistsUsername 检查用户名是否存在
	ExistsUsername(ctx context.Context, username string) (bool, error)

	// ExistsEmail 检查邮箱是否存在
	ExistsEmail(ctx context.Context, email string) (bool, error)
}
