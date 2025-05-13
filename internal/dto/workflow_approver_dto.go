package dto

import (
	"fmt"
	"slices"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
)

type ApproveRejectRequest struct {
	DocumentID   string                      `json:"document_id"`
	Level        int                         `json:"level"`
	Status       enums.AccountApprovalStatus `json:"status"`
	WorkflowType string                      `json:"workflow_type"`
	UserID       string
}

type ApproveRejectResponse struct {
	DocumentID    string `json:"document_id"`
	ApproveStatus string `json:"approve_status"`
}

func (w *ApproveRejectRequest) ValidationRequest() error {
	if w.DocumentID == "" {
		return fmt.Errorf("DocumentID is required")
	}

	if w.Level <= 0 {
		return fmt.Errorf("level must be greater than 0")
	}

	if w.WorkflowType == "" {
		return fmt.Errorf("workflow type is required")
	}

	if !slices.Contains([]enums.AccountApprovalStatus{
		enums.AccountApprovalStatusApproved,
		enums.AccountApprovalStatusRejected},
		w.Status,
	) {
		return fmt.Errorf("status cannot be other than Approved or Rejected")
	}

	if !slices.Contains([]string{
		common.WorkflowDepositApprover,
		common.WorkflowWithdrawalApprover},
		w.WorkflowType,
	) {
		return fmt.Errorf("workflow type cannot be other than deposit-approver or withdrawal-approver")
	}

	return nil
}

type ApiResponse struct {
	base.ApiResponse
}
