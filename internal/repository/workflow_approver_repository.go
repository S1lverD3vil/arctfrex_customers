package repository

import (
	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/model"
)

type WorkflowApproverRepository interface {
	CreateBulk(workflowApprover []model.WorkflowApprover) error
	GetWorkflowApproverByDocumentId(documentId string) ([]model.WorkflowApprover, error)
	UpdateApproverStatus(workflowApprover *model.WorkflowApprover) error
}

type workflowApproverRepository struct {
	db *gorm.DB
}

func NewWorkflowApproverRepository(db *gorm.DB) WorkflowApproverRepository {
	return &workflowApproverRepository{db: db}
}

func (dr *workflowApproverRepository) CreateBulk(workflowApprover []model.WorkflowApprover) error {
	return dr.db.Create(&workflowApprover).Error
}

func (dr *workflowApproverRepository) GetWorkflowApproverByDocumentId(documentId string) ([]model.WorkflowApprover, error) {
	var workflowApprover []model.WorkflowApprover
	queryParams := model.WorkflowApprover{
		DocumentID: documentId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Find(&workflowApprover, &queryParams).Error; err != nil {
		return nil, err
	}

	return workflowApprover, nil
}

func (dr *workflowApproverRepository) UpdateApproverStatus(workflowApprover *model.WorkflowApprover) error {
	return dr.db.Select(
		"status",
		"approved_at",
		"approved_by",
	).Updates(workflowApprover).Error
}
