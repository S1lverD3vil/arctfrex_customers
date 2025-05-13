package repository

import (
	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/model"
)

type WorkflowSettingRepository interface {
	GetWorkflowSettingByWorkflowType(workflowType string) (*model.WorkflowSetting, error)
}

type workflowSettingRepository struct {
	db *gorm.DB
}

func NewWorkflowSettingRepository(db *gorm.DB) WorkflowSettingRepository {
	return &workflowSettingRepository{db: db}
}

func (dr *workflowSettingRepository) GetWorkflowSettingByWorkflowType(workflowType string) (*model.WorkflowSetting, error) {
	var workflowSettingRepository model.WorkflowSetting
	queryParams := model.WorkflowSetting{
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
