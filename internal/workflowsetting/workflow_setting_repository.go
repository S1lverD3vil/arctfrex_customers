package workflowsetting

import (
	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
)

type workflowSettingRepository struct {
	db *gorm.DB
}

func NewWorkflowSettingRepository(db *gorm.DB) WorkflowSettingRepository {
	return &workflowSettingRepository{db: db}
}

func (dr *workflowSettingRepository) GetWorkflowSettingByWorkflowType(workflowType string) (*WorkflowSetting, error) {
	var workflowSettingRepository WorkflowSetting
	queryParams := WorkflowSetting{
		WorkflowType: workflowType,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Where(&queryParams).First(&workflowSettingRepository).Error; err != nil {
		return nil, err
	}

	return &workflowSettingRepository, nil
}
