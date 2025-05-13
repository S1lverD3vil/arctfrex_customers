package model

import (
	"time"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"
)

type WorkflowApprover struct {
	ID         string                      `json:"id" gorm:"primaryKey"`
	Level      int                         `json:"level"`
	DocumentID string                      `json:"document_id"`
	Status     enums.AccountApprovalStatus `json:"status"`
	ApprovedBy string                      `json:"approved_by"`
	ApprovedAt time.Time                   `json:"approved_at"`

	WorkflowSettingID *string `json:"workflow_setting_id"`

	WorkflowSetting *WorkflowSetting `gorm:"foreignKey:WorkflowSettingID"`

	base.BaseModel
}
