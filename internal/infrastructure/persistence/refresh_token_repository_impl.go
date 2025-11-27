package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/repository"
)

// RefreshTokenRepositoryImpl 刷新令牌仓储实现
type RefreshTokenRepositoryImpl struct {
	db *gorm.DB
}

// NewRefreshTokenRepository 创建刷新令牌仓储实例
func NewRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{db: db}
}

// Create 创建刷新令牌
func (r *RefreshTokenRepositoryImpl) Create(ctx context.Context, token *entity.RefreshToken) error {
	po := &RefreshTokenPO{}
	po.FromEntity(token)

	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}

	// 回填ID
	token.ID = po.ID
	token.CreatedAt = po.CreatedAt
	token.UpdatedAt = po.UpdatedAt

	return nil
}

// GetByToken 根据token字符串获取令牌
func (r *RefreshTokenRepositoryImpl) GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var po RefreshTokenPO
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return po.ToEntity(), nil
}

// Revoke 撤销令牌
func (r *RefreshTokenRepositoryImpl) Revoke(ctx context.Context, token string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&RefreshTokenPO{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

// RevokeAllByUserID 撤销用户的所有令牌（用户登出所有设备）
func (r *RefreshTokenRepositoryImpl) RevokeAllByUserID(ctx context.Context, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&RefreshTokenPO{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

// DeleteExpired 删除过期的令牌（定期清理）
func (r *RefreshTokenRepositoryImpl) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&RefreshTokenPO{}).Error
}

// GetValidTokensByUserID 获取用户的有效令牌列表
func (r *RefreshTokenRepositoryImpl) GetValidTokensByUserID(ctx context.Context, userID int64) ([]*entity.RefreshToken, error) {
	var pos []*RefreshTokenPO
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Find(&pos).Error; err != nil {
		return nil, err
	}

	tokens := make([]*entity.RefreshToken, len(pos))
	for i, po := range pos {
		tokens[i] = po.ToEntity()
	}

	return tokens, nil
}
