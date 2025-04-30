package workflowsetting

import (
	"arctfrex-customers/internal/base"
)

type WorkflowSetting struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Config       string `json:"config" gorm:"type:jsonb"`
	WorkflowType string `json:"workflow_type"`

	base.BaseModel
}
