package entity

import "time"

// RefreshToken 刷新令牌实体
type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	DeviceID  string
	UserAgent string
	IPAddress string
	IsRevoked bool
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsExpired 判断令牌是否已过期
func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsValid 判断令牌是否有效
func (r *RefreshToken) IsValid() bool {
	return !r.IsRevoked && !r.IsExpired()
}
