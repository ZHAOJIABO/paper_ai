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
