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
