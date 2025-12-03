package persistence

import (
	"encoding/json"
	"time"

	"paper_ai/internal/domain/entity"
	"gorm.io/gorm"
)

// PolishRecordPO 润色记录持久化对象（PO: Persistent Object）
// 包含ORM框架的标签，与entity分离，保持领域层纯净
type PolishRecordPO struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	TraceID         string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_trace_id;comment:'润色追踪ID(纯数字)'"`
	UserID          int64          `gorm:"not null;index:idx_user_id"` // 用户ID

	OriginalContent string         `gorm:"type:text;not null"`
	Style           string         `gorm:"type:varchar(20);not null;index:idx_style"`
	Language        string         `gorm:"type:varchar(10);not null;index:idx_language"`

	PolishedContent string         `gorm:"type:text;not null"`
	OriginalLength  int            `gorm:"not null"`
	PolishedLength  int            `gorm:"not null"`

	Provider        string         `gorm:"type:varchar(50);not null;index:idx_provider"`
	Model           string         `gorm:"type:varchar(100);not null"`

	Mode            string         `gorm:"type:varchar(20);not null;default:'single';index:idx_mode;comment:'润色模式: single(单版本) / multi(多版本)'"`

	ProcessTimeMs   int            `gorm:"default:0;index:idx_process_time"`

	Status          string         `gorm:"type:varchar(20);not null;default:'success';index:idx_status"`
	ErrorMessage    string         `gorm:"type:text"`

	// 对比数据（新增）
	ComparisonData  *string `gorm:"type:jsonb"` // 存储对比数据的JSON（指针类型，允许NULL）
	ChangesCount    int     `gorm:"default:0"`  // 修改总数
	AcceptedChanges *string `gorm:"type:jsonb"` // 用户接受的修改ID列表（指针类型，允许NULL）
	RejectedChanges *string `gorm:"type:jsonb"` // 用户拒绝的修改ID列表（指针类型，允许NULL）
	FinalContent    string  `gorm:"type:text"`  // 用户最终确认的文本（原文+接受的修改）

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
	record := &entity.PolishRecord{
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
		Mode:            po.Mode,
		ProcessTimeMs:   po.ProcessTimeMs,
		Status:          po.Status,
		ErrorMessage:    po.ErrorMessage,
		ComparisonData:  func() string {
			if po.ComparisonData != nil {
				return *po.ComparisonData
			}
			return ""
		}(),
		ChangesCount:    po.ChangesCount,
		FinalContent:    po.FinalContent,
		CreatedAt:       po.CreatedAt,
		UpdatedAt:       po.UpdatedAt,
	}

	// 解析 JSON 数组
	if po.AcceptedChanges != nil && *po.AcceptedChanges != "" {
		var acceptedIDs []string
		if err := json.Unmarshal([]byte(*po.AcceptedChanges), &acceptedIDs); err == nil {
			record.AcceptedChanges = acceptedIDs
		} else {
			record.AcceptedChanges = []string{}
		}
	} else {
		record.AcceptedChanges = []string{}
	}

	if po.RejectedChanges != nil && *po.RejectedChanges != "" {
		var rejectedIDs []string
		if err := json.Unmarshal([]byte(*po.RejectedChanges), &rejectedIDs); err == nil {
			record.RejectedChanges = rejectedIDs
		} else {
			record.RejectedChanges = []string{}
		}
	} else {
		record.RejectedChanges = []string{}
	}

	return record
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
	po.Mode = e.Mode
	po.ProcessTimeMs = e.ProcessTimeMs
	po.Status = e.Status
	po.ErrorMessage = e.ErrorMessage

	// 处理 ComparisonData JSONB 字段（空字符串设为 nil）
	if e.ComparisonData != "" {
		po.ComparisonData = &e.ComparisonData
	} else {
		po.ComparisonData = nil
	}

	po.ChangesCount = e.ChangesCount
	po.FinalContent = e.FinalContent
	po.CreatedAt = e.CreatedAt
	po.UpdatedAt = e.UpdatedAt

	// 序列化 JSON 数组（如果为空则设为 nil，让数据库存储 NULL）
	if len(e.AcceptedChanges) > 0 {
		jsonBytes, err := json.Marshal(e.AcceptedChanges)
		if err == nil {
			jsonStr := string(jsonBytes)
			po.AcceptedChanges = &jsonStr
		} else {
			po.AcceptedChanges = nil
		}
	} else {
		po.AcceptedChanges = nil
	}

	if len(e.RejectedChanges) > 0 {
		jsonBytes, err := json.Marshal(e.RejectedChanges)
		if err == nil {
			jsonStr := string(jsonBytes)
			po.RejectedChanges = &jsonStr
		} else {
			po.RejectedChanges = nil
		}
	} else {
		po.RejectedChanges = nil
	}
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

	// 多版本润色功能权限
	EnableMultiVersion bool `gorm:"default:false;index:idx_enable_multi_version;comment:'是否启用多版本功能'"`
	MultiVersionQuota  int  `gorm:"default:0;comment:'多版本配额(0=无限)'"`

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
		EnableMultiVersion: po.EnableMultiVersion,
		MultiVersionQuota: po.MultiVersionQuota,
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
	po.EnableMultiVersion = e.EnableMultiVersion
	po.MultiVersionQuota = e.MultiVersionQuota
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

// PolishVersionPO 润色版本持久化对象
type PolishVersionPO struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	RecordID int64  `gorm:"not null;index:idx_record_id"`

	// 版本信息
	VersionType string `gorm:"type:varchar(32);not null;index:idx_version_type"`

	// 输出内容
	PolishedContent string  `gorm:"type:text;not null"`
	PolishedLength  int     `gorm:"not null"`
	Suggestions     *string `gorm:"type:jsonb"` // JSON数组

	// AI信息
	ModelUsed string `gorm:"type:varchar(64);not null"`
	PromptID  int64  `gorm:""`

	// 性能指标
	ProcessTimeMs int `gorm:"not null;default:0"`

	// 状态
	Status       string `gorm:"type:varchar(20);not null;default:'success'"`
	ErrorMessage string `gorm:"type:text"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 指定表名
func (PolishVersionPO) TableName() string {
	return "polish_versions"
}

// ToEntity 转换为领域实体
func (po *PolishVersionPO) ToEntity() *entity.PolishVersion {
	version := &entity.PolishVersion{
		ID:              po.ID,
		RecordID:        po.RecordID,
		VersionType:     po.VersionType,
		PolishedContent: po.PolishedContent,
		PolishedLength:  po.PolishedLength,
		ModelUsed:       po.ModelUsed,
		PromptID:        po.PromptID,
		ProcessTimeMs:   po.ProcessTimeMs,
		Status:          po.Status,
		ErrorMessage:    po.ErrorMessage,
		CreatedAt:       po.CreatedAt,
	}

	// 解析 JSON 数组
	if po.Suggestions != nil && *po.Suggestions != "" {
		var suggestions []string
		if err := json.Unmarshal([]byte(*po.Suggestions), &suggestions); err == nil {
			version.Suggestions = suggestions
		} else {
			version.Suggestions = []string{}
		}
	} else {
		version.Suggestions = []string{}
	}

	return version
}

// FromEntity 从领域实体创建PO
func (po *PolishVersionPO) FromEntity(e *entity.PolishVersion) {
	po.ID = e.ID
	po.RecordID = e.RecordID
	po.VersionType = e.VersionType
	po.PolishedContent = e.PolishedContent
	po.PolishedLength = e.PolishedLength
	po.ModelUsed = e.ModelUsed
	po.PromptID = e.PromptID
	po.ProcessTimeMs = e.ProcessTimeMs
	po.Status = e.Status
	po.ErrorMessage = e.ErrorMessage
	po.CreatedAt = e.CreatedAt

	// 序列化 JSON 数组
	if len(e.Suggestions) > 0 {
		jsonBytes, err := json.Marshal(e.Suggestions)
		if err == nil {
			jsonStr := string(jsonBytes)
			po.Suggestions = &jsonStr
		} else {
			po.Suggestions = nil
		}
	} else {
		po.Suggestions = nil
	}
}

// PolishPromptPO Prompt模板持久化对象
type PolishPromptPO struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	// 基本信息
	Name        string `gorm:"type:varchar(128);not null"`
	VersionType string `gorm:"type:varchar(32);not null;index:idx_version_type_prompt"`
	Language    string `gorm:"type:varchar(16);not null;index:idx_language_prompt"`
	Style       string `gorm:"type:varchar(32);not null;index:idx_style_prompt"`

	// Prompt内容
	SystemPrompt       string `gorm:"type:text;not null"`
	UserPromptTemplate string `gorm:"type:text;not null"`

	// 版本管理
	Version  int  `gorm:"not null;default:1"`
	IsActive bool `gorm:"default:true;index:idx_active_prompt"`

	// 元数据
	Description string  `gorm:"type:text"`
	Tags        *string `gorm:"type:jsonb"`

	// A/B测试
	ABTestGroup string `gorm:"type:varchar(32)"`
	Weight      int    `gorm:"default:100"`

	// 统计信息
	UsageCount      int     `gorm:"default:0"`
	SuccessRate     float64 `gorm:"type:decimal(5,2)"`
	AvgSatisfaction float64 `gorm:"type:decimal(3,2)"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_created_at_prompt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy string    `gorm:"type:varchar(128)"`
}

// TableName 指定表名
func (PolishPromptPO) TableName() string {
	return "polish_prompts"
}

// ToEntity 转换为领域实体
func (po *PolishPromptPO) ToEntity() *entity.PolishPrompt {
	prompt := &entity.PolishPrompt{
		ID:                 po.ID,
		Name:               po.Name,
		VersionType:        po.VersionType,
		Language:           po.Language,
		Style:              po.Style,
		SystemPrompt:       po.SystemPrompt,
		UserPromptTemplate: po.UserPromptTemplate,
		Version:            po.Version,
		IsActive:           po.IsActive,
		Description:        po.Description,
		ABTestGroup:        po.ABTestGroup,
		Weight:             po.Weight,
		UsageCount:         po.UsageCount,
		SuccessRate:        po.SuccessRate,
		AvgSatisfaction:    po.AvgSatisfaction,
		CreatedAt:          po.CreatedAt,
		UpdatedAt:          po.UpdatedAt,
		CreatedBy:          po.CreatedBy,
	}

	// 解析 JSON 数组
	if po.Tags != nil && *po.Tags != "" {
		var tags []string
		if err := json.Unmarshal([]byte(*po.Tags), &tags); err == nil {
			prompt.Tags = tags
		} else {
			prompt.Tags = []string{}
		}
	} else {
		prompt.Tags = []string{}
	}

	return prompt
}

// FromEntity 从领域实体创建PO
func (po *PolishPromptPO) FromEntity(e *entity.PolishPrompt) {
	po.ID = e.ID
	po.Name = e.Name
	po.VersionType = e.VersionType
	po.Language = e.Language
	po.Style = e.Style
	po.SystemPrompt = e.SystemPrompt
	po.UserPromptTemplate = e.UserPromptTemplate
	po.Version = e.Version
	po.IsActive = e.IsActive
	po.Description = e.Description
	po.ABTestGroup = e.ABTestGroup
	po.Weight = e.Weight
	po.UsageCount = e.UsageCount
	po.SuccessRate = e.SuccessRate
	po.AvgSatisfaction = e.AvgSatisfaction
	po.CreatedAt = e.CreatedAt
	po.UpdatedAt = e.UpdatedAt
	po.CreatedBy = e.CreatedBy

	// 序列化 JSON 数组
	if len(e.Tags) > 0 {
		jsonBytes, err := json.Marshal(e.Tags)
		if err == nil {
			jsonStr := string(jsonBytes)
			po.Tags = &jsonStr
		} else {
			po.Tags = nil
		}
	} else {
		po.Tags = nil
	}
}
