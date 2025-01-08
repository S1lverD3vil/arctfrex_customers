package withdrawal

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"

	"gorm.io/gorm"
)

type withdrawalRepository struct {
	db *gorm.DB
}

func NewWithdrawalRepository(db *gorm.DB) WithdrawalRepository {
	return &withdrawalRepository{db: db}
}

func (dr *withdrawalRepository) Create(withdrawal *Withdrawal) error {
	return dr.db.Create(withdrawal).Error
}

func (dr *withdrawalRepository) GetPendingAccountByAccountIdUserId(accountId, userId string) (*Withdrawal, error) {
	var withdrawal Withdrawal
	queryParams := Withdrawal{
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

func (dr *withdrawalRepository) GetWithdrawalsByUserIdAccountId(userId, accountId string) (*[]Withdrawals, error) {
	var withdrawals []Withdrawals
	if err := dr.db.Table("withdrawals").
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Select("withdrawals.id as withdrawalid, withdrawals.account_id as accountid, withdrawals.user_id as userid, users.name as name, users.email, withdrawals.amount, withdrawals.approval_status as approval_status, withdrawals.created_at as transaction_date").
		Where("withdrawals.approval_status != ? AND withdrawals.is_active = ? AND withdrawals.user_id = ? AND withdrawals.account_id = ?", 0, true, userId, accountId).
		Scan(&withdrawals).Error; err != nil {

		return nil, err
	}

	return &withdrawals, nil
}

func (dr *withdrawalRepository) GetWithdrawalByIdUserId(userId, withdrawalId string) (*Withdrawal, error) {
	var withdrawal Withdrawal
	queryParams := Withdrawal{
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

func (dr *withdrawalRepository) GetBackOfficePendingWithdrawals() (*[]BackOfficePendingWithdrawal, error) {
	var backOfficePendingWithdrawals []BackOfficePendingWithdrawal
	if err := dr.db.Table("withdrawals").
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Select("withdrawals.id as withdrawalid, withdrawals.account_id as accountid, withdrawals.user_id as userid, users.name as name, users.email, withdrawals.amount, withdrawals.approval_status as approval_status").
		Where("withdrawals.approval_status = ? AND withdrawals.is_active = ?", enums.WithdrawalApprovalStatusPending, true).
		Scan(&backOfficePendingWithdrawals).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingWithdrawals, nil
}

func (dr *withdrawalRepository) GetBackOfficePendingWithdrawalDetail(withdrawalId string) (*BackOfficePendingWithdrawalDetail, error) {
	var backOfficePendingWithdrawalDetail BackOfficePendingWithdrawalDetail
	if err := dr.db.Table("withdrawals").
		Joins("JOIN users ON users.id = withdrawals.user_id").
		Select("withdrawals.id as withdrawalid, withdrawals.account_id as accountid, withdrawals.user_id as userid, users.name as name, users.email, withdrawals.amount, withdrawals.bank_name, withdrawals.bank_account_number, withdrawals.bank_beneficiary_name, withdrawals.approval_status as approval_status").
		Where("withdrawals.approval_status = ? AND withdrawals.is_active = ?", enums.WithdrawalApprovalStatusPending, true).
		Scan(&backOfficePendingWithdrawalDetail).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingWithdrawalDetail, nil
}

func (dr *withdrawalRepository) GetPendingWithdrawalsById(withdrawalId string) (*Withdrawal, error) {
	var withdrawals Withdrawal
	queryParams := Withdrawal{
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

func (dr *withdrawalRepository) UpdateWithdrawalApprovalStatus(withdrawal *Withdrawal) error {
	return dr.db.Select(
		"ApprovalStatus",
		"ApprovedBy",
		"ApprovedAt",
	).Updates(withdrawal).Error
}
