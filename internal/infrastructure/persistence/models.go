package persistence

import (
	"time"

	"paper_ai/internal/domain/entity"
	"gorm.io/gorm"
)

// PolishRecordPO 润色记录持久化对象（PO: Persistent Object）
// 包含ORM框架的标签，与entity分离，保持领域层纯净
type PolishRecordPO struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	TraceID         string         `gorm:"type:varchar(64);not null;uniqueIndex:idx_trace_id"`
	UserID          int64          `gorm:"not null;index:idx_user_id"` // 用户ID

	OriginalContent string         `gorm:"type:text;not null"`
	Style           string         `gorm:"type:varchar(20);not null;index:idx_style"`
	Language        string         `gorm:"type:varchar(10);not null;index:idx_language"`

	PolishedContent string         `gorm:"type:text;not null"`
	OriginalLength  int            `gorm:"not null"`
	PolishedLength  int            `gorm:"not null"`

	Provider        string         `gorm:"type:varchar(50);not null;index:idx_provider"`
	Model           string         `gorm:"type:varchar(100);not null"`

	ProcessTimeMs   int            `gorm:"default:0;index:idx_process_time"`

	Status          string         `gorm:"type:varchar(20);not null;default:'success';index:idx_status"`
	ErrorMessage    string         `gorm:"type:text"`

	CreatedAt       time.Time      `gorm:"autoCreateTime;index:idx_created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"` // 软删除支持
}

// TableName 指定表名
func (PolishRecordPO) TableName() string {
	return "polish_records"
}

// ToEntity 转换为领域实体
func (po *PolishRecordPO) ToEntity() *entity.PolishRecord {
	return &entity.PolishRecord{
		ID:              po.ID,
		TraceID:         po.TraceID,
		UserID:          po.UserID,
		OriginalContent: po.OriginalContent,
		Style:           po.Style,
		Language:        po.Language,
		PolishedContent: po.PolishedContent,
		OriginalLength:  po.OriginalLength,
		PolishedLength:  po.PolishedLength,
		Provider:        po.Provider,
		Model:           po.Model,
		ProcessTimeMs:   po.ProcessTimeMs,
		Status:          po.Status,
		ErrorMessage:    po.ErrorMessage,
		CreatedAt:       po.CreatedAt,
		UpdatedAt:       po.UpdatedAt,
	}
}

// FromEntity 从领域实体创建PO
func (po *PolishRecordPO) FromEntity(e *entity.PolishRecord) {
	po.ID = e.ID
	po.TraceID = e.TraceID
	po.UserID = e.UserID
	po.OriginalContent = e.OriginalContent
	po.Style = e.Style
	po.Language = e.Language
	po.PolishedContent = e.PolishedContent
	po.OriginalLength = e.OriginalLength
	po.PolishedLength = e.PolishedLength
	po.Provider = e.Provider
	po.Model = e.Model
	po.ProcessTimeMs = e.ProcessTimeMs
	po.Status = e.Status
	po.ErrorMessage = e.ErrorMessage
	po.CreatedAt = e.CreatedAt
	po.UpdatedAt = e.UpdatedAt
}

// ToEntityList 批量转换为实体列表
func ToEntityList(pos []*PolishRecordPO) []*entity.PolishRecord {
	entities := make([]*entity.PolishRecord, len(pos))
	for i, po := range pos {
		entities[i] = po.ToEntity()
	}
	return entities
}

// UserPO 用户持久化对象
type UserPO struct {
	ID               int64          `gorm:"primaryKey"`  // 使用Snowflake ID，不再自增
	Username         string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_username"`
	Email            string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_email"`
	PasswordHash     string         `gorm:"type:varchar(255);not null"`
	Nickname         string         `gorm:"type:varchar(50)"`
	AvatarURL        string         `gorm:"type:varchar(255)"`
	Status           string         `gorm:"type:varchar(20);not null;default:'active';index:idx_status"`
	EmailVerified    bool           `gorm:"default:false"`
	LastLoginAt      *time.Time     `gorm:""`
	LastLoginIP      string         `gorm:"type:varchar(50)"`
	LoginCount       int            `gorm:"default:0"`
	FailedLoginCount int            `gorm:"default:0"`
	CreatedAt        time.Time      `gorm:"autoCreateTime;index:idx_created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

// TableName 指定表名
func (UserPO) TableName() string {
	return "users"
}

// ToEntity 转换为领域实体
func (po *UserPO) ToEntity() *entity.User {
	return &entity.User{
		ID:               po.ID,
		Username:         po.Username,
		Email:            po.Email,
		PasswordHash:     po.PasswordHash,
		Nickname:         po.Nickname,
		AvatarURL:        po.AvatarURL,
		Status:           po.Status,
		EmailVerified:    po.EmailVerified,
		LastLoginAt:      po.LastLoginAt,
		LastLoginIP:      po.LastLoginIP,
		LoginCount:       po.LoginCount,
		FailedLoginCount: po.FailedLoginCount,
		CreatedAt:        po.CreatedAt,
		UpdatedAt:        po.UpdatedAt,
	}
}

// FromEntity 从领域实体创建PO
func (po *UserPO) FromEntity(e *entity.User) {
	po.ID = e.ID
	po.Username = e.Username
	po.Email = e.Email
	po.PasswordHash = e.PasswordHash
	po.Nickname = e.Nickname
	po.AvatarURL = e.AvatarURL
	po.Status = e.Status
	po.EmailVerified = e.EmailVerified
	po.LastLoginAt = e.LastLoginAt
	po.LastLoginIP = e.LastLoginIP
	po.LoginCount = e.LoginCount
	po.FailedLoginCount = e.FailedLoginCount
	po.CreatedAt = e.CreatedAt
	po.UpdatedAt = e.UpdatedAt
}

// RefreshTokenPO 刷新令牌持久化对象
type RefreshTokenPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"not null;index:idx_user_id"`
	Token     string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_token"`
	ExpiresAt time.Time      `gorm:"not null;index:idx_expires_at"`
	DeviceID  string         `gorm:"type:varchar(100)"`
	UserAgent string         `gorm:"type:varchar(500)"`
	IPAddress string         `gorm:"type:varchar(50)"`
	IsRevoked bool           `gorm:"default:false"`
	RevokedAt *time.Time     `gorm:""`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定表名
func (RefreshTokenPO) TableName() string {
	return "refresh_tokens"
}

// ToEntity 转换为领域实体
func (po *RefreshTokenPO) ToEntity() *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        po.ID,
		UserID:    po.UserID,
		Token:     po.Token,
		ExpiresAt: po.ExpiresAt,
		DeviceID:  po.DeviceID,
		UserAgent: po.UserAgent,
		IPAddress: po.IPAddress,
		IsRevoked: po.IsRevoked,
		RevokedAt: po.RevokedAt,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

// FromEntity 从领域实体创建PO
func (po *RefreshTokenPO) FromEntity(e *entity.RefreshToken) {
	po.ID = e.ID
	po.UserID = e.UserID
	po.Token = e.Token
	po.ExpiresAt = e.ExpiresAt
	po.DeviceID = e.DeviceID
	po.UserAgent = e.UserAgent
	po.IPAddress = e.IPAddress
	po.IsRevoked = e.IsRevoked
	po.RevokedAt = e.RevokedAt
	po.CreatedAt = e.CreatedAt
	po.UpdatedAt = e.UpdatedAt
}
