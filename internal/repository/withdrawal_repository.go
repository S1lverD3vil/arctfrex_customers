package repository

import (
	"fmt"

	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
)

type WithdrawalRepository interface {
	Create(withdrawal *model.Withdrawal) error
	GetPendingAccountByAccountIdUserId(accountId, userId string) (*model.Withdrawal, error)
	GetWithdrawalsByUserIdAccountId(userId, accountId string) (*[]model.Withdrawals, error)
	GetWithdrawalByIdUserId(userId, withdrawalId string) (*model.Withdrawal, error)
	GetBackOfficePendingWithdrawals() (*[]model.BackOfficePendingWithdrawal, error)
	GetBackOfficePendingWithdrawalSPA(request model.WithdrawalBackOfficeParam) ([]model.BackOfficePendingWithdrawal, error)
	GetBackOfficePendingWithdrawalMulti(request model.WithdrawalBackOfficeParam) ([]model.BackOfficePendingWithdrawal, error)
	GetBackOfficePendingWithdrawalDetail(withdrawalId string) (*model.BackOfficePendingWithdrawalDetail, error)
	GetPendingWithdrawalsById(withdrawalId string) (*model.Withdrawal, error)
	UpdateWithdrawalApprovalStatus(withdrawal *model.Withdrawal) error
}

type withdrawalRepository struct {
	db *gorm.DB
}

func NewWithdrawalRepository(db *gorm.DB) WithdrawalRepository {
	return &withdrawalRepository{db: db}
}

func (dr *withdrawalRepository) Create(withdrawal *model.Withdrawal) error {
	return dr.db.Create(withdrawal).Error
}

func (dr *withdrawalRepository) GetPendingAccountByAccountIdUserId(accountId, userId string) (*model.Withdrawal, error) {
	var withdrawal model.Withdrawal
	queryParams := model.Withdrawal{
		AccountID:      accountId,
		UserID:         userId,
		ApprovalStatus: enums.WithdrawalApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Where(&queryParams).First(&withdrawal).Error; err != nil {
		return nil, err
	}

	return &withdrawal, nil
}

func (dr *withdrawalRepository) GetWithdrawalsByUserIdAccountId(userId, accountId string) (*[]model.Withdrawals, error) {
	var withdrawals []model.Withdrawals
	if err := dr.db.Table("withdrawals").
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Select("withdrawals.id as withdrawalid, withdrawals.account_id as accountid, withdrawals.user_id as userid, users.name as name, users.email, withdrawals.amount, withdrawals.approval_status as approval_status, withdrawals.created_at as transaction_date").
		Where("withdrawals.approval_status != ? AND withdrawals.is_active = ? AND withdrawals.user_id = ? AND withdrawals.account_id = ?", 0, true, userId, accountId).
		Scan(&withdrawals).Error; err != nil {

		return nil, err
	}

	return &withdrawals, nil
}

func (dr *withdrawalRepository) GetWithdrawalByIdUserId(userId, withdrawalId string) (*model.Withdrawal, error) {
	var withdrawal model.Withdrawal
	queryParams := model.Withdrawal{
		ID:     withdrawalId,
		UserID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Where(&queryParams).First(&withdrawal).Error; err != nil {
		return nil, err
	}

	return &withdrawal, nil
}

func (dr *withdrawalRepository) GetBackOfficePendingWithdrawals() (*[]model.BackOfficePendingWithdrawal, error) {
	var backOfficePendingWithdrawals []model.BackOfficePendingWithdrawal
	if err := dr.db.Table("withdrawals").
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Select("withdrawals.id as withdrawalid, withdrawals.account_id as accountid, withdrawals.user_id as userid, users.name as name, users.email, withdrawals.amount, withdrawals.approval_status as approval_status").
		Where("withdrawals.approval_status = ? AND withdrawals.is_active = ?", enums.WithdrawalApprovalStatusPending, true).
		Scan(&backOfficePendingWithdrawals).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingWithdrawals, nil
}

func (dr *withdrawalRepository) GetBackOfficePendingWithdrawalDetail(withdrawalId string) (*model.BackOfficePendingWithdrawalDetail, error) {
	var backOfficePendingWithdrawalDetail model.BackOfficePendingWithdrawalDetail
	if err := dr.db.Table("withdrawals").
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Select("withdrawals.id as withdrawalid, withdrawals.account_id as accountid, withdrawals.user_id as userid, users.name as name, users.email, withdrawals.amount, withdrawals.bank_name, withdrawals.bank_account_number, withdrawals.bank_beneficiary_name, withdrawals.approval_status as approval_status").
		Where("withdrawals.approval_status = ? AND withdrawals.is_active = ?", enums.WithdrawalApprovalStatusPending, true).
		Scan(&backOfficePendingWithdrawalDetail).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingWithdrawalDetail, nil
}

func (dr *withdrawalRepository) GetPendingWithdrawalsById(withdrawalId string) (*model.Withdrawal, error) {
	var withdrawals model.Withdrawal
	queryParams := model.Withdrawal{
		ID:             withdrawalId,
		ApprovalStatus: enums.WithdrawalApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Find(&withdrawals, &queryParams).Error; err != nil {
		return nil, err
	}

	return &withdrawals, nil
}

func (dr *withdrawalRepository) UpdateWithdrawalApprovalStatus(withdrawal *model.Withdrawal) error {
	return dr.db.Select(
		"ApprovalStatus",
		"ApprovedBy",
		"ApprovedAt",
	).Updates(withdrawal).Error
}

func (dr *withdrawalRepository) GetBackOfficePendingWithdrawalSPA(request model.WithdrawalBackOfficeParam) (backOfficePendingWithdrawals []model.BackOfficePendingWithdrawal, err error) {
	baseSelect := `
		withdrawals.id AS withdrawal_id,
		withdrawals.account_id,
		withdrawals.user_id, 
		users.name AS name,
		users.email,
		users.meta_login_id, 
		withdrawals.amount, 
		withdrawals.amount_usd,
		withdrawals.approval_status,
		withdrawals.created_at,
		withdrawals.bank_name
	`

	if request.Menutype == common.Finance {
		baseSelect += `,
		bu.name AS settlement_by`
	}

	query := dr.db.Table("withdrawals").
		Select(baseSelect).
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Where("withdrawals.approval_status = ? AND withdrawals.is_active = ?", enums.DepositApprovalStatusPending, true)

	switch request.Menutype {
	case common.Settlement:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = withdrawals.id").
			Where("wa1.level=1 AND wa1.is_active=? AND wa1.status=?", true, enums.AccountApprovalStatusPending)
	case common.Finance:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = withdrawals.id AND wa1.level = 1 AND wa1.status = ? AND wa1.is_active = ?", enums.AccountApprovalStatusApproved, true).
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = withdrawals.id").
			Where("wa2.level=2 AND wa2.is_active=? AND wa2.status=?", true, enums.AccountApprovalStatusPending)
	default:
		return backOfficePendingWithdrawals, fmt.Errorf("invalid menu type: %s", request.Menutype)
	}

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize

	if err = query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err = query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficePendingWithdrawals).Error; err != nil {
		return nil, err
	}

	return backOfficePendingWithdrawals, nil
}

func (dr *withdrawalRepository) GetBackOfficePendingWithdrawalMulti(request model.WithdrawalBackOfficeParam) (backOfficePendingWithdrawals []model.BackOfficePendingWithdrawal, err error) {
	baseSelect := `
		withdrawals.id AS withdrawal_id,
		withdrawals.account_id,
		withdrawals.user_id, 
		users.name AS name,
		users.email,
		users.meta_login_id, 
		withdrawals.amount, 
		withdrawals.amount_usd,
		withdrawals.approval_status,
		withdrawals.created_at,
		withdrawals.bank_name
	`

	if request.Menutype == common.Finance {
		baseSelect += `,
		bu.name AS settlement_by`
	}

	query := dr.db.Table("withdrawals").
		Select(baseSelect).
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Where("withdrawals.approval_status = ? AND withdrawals.is_active = ?", enums.DepositApprovalStatusPending, true)

	switch request.Menutype {
	case common.Settlement:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = withdrawals.id").
			Where("wa1.level=1 AND wa1.is_active=? AND wa1.status=?", true, enums.AccountApprovalStatusPending)
	case common.Finance:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = withdrawals.id AND wa1.level = 1 AND wa1.status = ? AND wa1.is_active = ?", enums.AccountApprovalStatusApproved, true).
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = withdrawals.id").
			Where("wa2.level=2 AND wa2.is_active=? AND wa2.status=?", true, enums.AccountApprovalStatusPending)
	default:
		return backOfficePendingWithdrawals, fmt.Errorf("invalid menu type: %s", request.Menutype)
	}

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize

	if err = query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err = query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficePendingWithdrawals).Error; err != nil {
		return nil, err
	}

	return backOfficePendingWithdrawals, nil
}
