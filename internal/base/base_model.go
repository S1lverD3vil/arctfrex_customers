package base

import (
	"time"

	"gorm.io/gorm"
)

func (b *BaseModel) BeforeCreate(db *gorm.DB) (err error) {
	if b.CreatedBy == "" {
		b.CreatedBy = "SYSTEM"
	}
	return
}

func (b *BaseModel) BeforeSave(db *gorm.DB) (err error) {
	if b.CreatedBy == "" {
		b.CreatedBy = "SYSTEM"
	}
	return
}

// BaseModel includes common fields for tracking creation and modification
type BaseModel struct {
	IsActive   bool      `json:"is_active"`
	CreatedBy  string    `gorm:"column:created_by;type:varchar(1000);"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"` // Auto-set on creation
	ModifiedBy string    `gorm:"column:modified_by;type:varchar(1000);"`
	ModifiedAt time.Time `gorm:"column:modified_at;autoUpdateTime"` // Auto-set on update
}
