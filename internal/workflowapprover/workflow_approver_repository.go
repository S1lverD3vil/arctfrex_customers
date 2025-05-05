package workflowapprover

import (
	"gorm.io/gorm"
)

type workflowApproverRepository struct {
	db *gorm.DB
}

func NewWorkflowApproverRepository(db *gorm.DB) WorkflowApproverRepository {
	return &workflowApproverRepository{db: db}
}

func (dr *workflowApproverRepository) CreateBulk(workflowApprover []WorkflowApprover) error {
	return dr.db.Create(&workflowApprover).Error
}
