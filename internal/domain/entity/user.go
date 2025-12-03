package entity

import "time"

// User 用户实体
type User struct {
	ID               int64
	Username         string
	Email            string
	PasswordHash     string
	Nickname         string
	AvatarURL        string
	Status           string // active, inactive, banned
	EmailVerified    bool
	LastLoginAt      *time.Time
	LastLoginIP      string
	LoginCount       int
	FailedLoginCount int

	// 多版本润色功能权限
	EnableMultiVersion  bool // 是否启用多版本功能
	MultiVersionQuota   int  // 多版本配额（0=无限）

	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// IsActive 判断用户是否处于活跃状态
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsBanned 判断用户是否被封禁
func (u *User) IsBanned() bool {
	return u.Status == "banned"
}

// HasMultiVersionPermission 判断用户是否有多版本权限
func (u *User) HasMultiVersionPermission() bool {
	return u.EnableMultiVersion
}

// HasUnlimitedQuota 判断用户是否有无限配额
func (u *User) HasUnlimitedQuota() bool {
	return u.MultiVersionQuota == 0
}
