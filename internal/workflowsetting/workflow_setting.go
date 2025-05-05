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

type WorkflowConfig struct {
	Approvers []Approver `json:"approvers"`
}

type Approver struct {
	Level  int    `json:"level"`
	RoleID string `json:"role_id"`
}

type WorkflowSettingRepository interface {
	GetWorkflowSettingByWorkflowType(workflowType string) (*WorkflowSetting, error)
}
