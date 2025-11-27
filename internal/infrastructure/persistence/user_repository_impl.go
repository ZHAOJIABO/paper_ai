package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/repository"
	"paper_ai/pkg/idgen"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	// 生成Snowflake ID
	id, err := idgen.GenerateID()
	if err != nil {
		return err
	}
	user.ID = id

	po := &UserPO{}
	po.FromEntity(user)

	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}

	// 回填时间戳
	user.CreatedAt = po.CreatedAt
	user.UpdatedAt = po.UpdatedAt

	return nil
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var po UserPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return po.ToEntity(), nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var po UserPO
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return po.ToEntity(), nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var po UserPO
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return po.ToEntity(), nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	po := &UserPO{}
	po.FromEntity(user)

	return r.db.WithContext(ctx).Save(po).Error
}

// UpdateLoginInfo 更新登录信息
func (r *UserRepositoryImpl) UpdateLoginInfo(ctx context.Context, userID int64, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&UserPO{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at":      now,
			"last_login_ip":      ip,
			"login_count":        gorm.Expr("login_count + 1"),
			"failed_login_count": 0, // 登录成功后重置失败次数
		}).Error
}

// IncrementFailedLoginCount 增加失败登录次数
func (r *UserRepositoryImpl) IncrementFailedLoginCount(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&UserPO{}).
		Where("id = ?", userID).
		Update("failed_login_count", gorm.Expr("failed_login_count + 1")).Error
}

// ResetFailedLoginCount 重置失败登录次数
func (r *UserRepositoryImpl) ResetFailedLoginCount(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&UserPO{}).
		Where("id = ?", userID).
		Update("failed_login_count", 0).Error
}

// ExistsUsername 检查用户名是否存在
func (r *UserRepositoryImpl) ExistsUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserPO{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsEmail 检查邮箱是否存在
func (r *UserRepositoryImpl) ExistsEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserPO{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
