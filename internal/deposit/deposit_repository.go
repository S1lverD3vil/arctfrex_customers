package deposit

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"

	"gorm.io/gorm"
)

type depositRepository struct {
	db *gorm.DB
}

func NewDepositRepository(db *gorm.DB) DepositRepository {
	return &depositRepository{db: db}
}

func (dr *depositRepository) Create(deposit *Deposit) error {
	return dr.db.Create(deposit).Error
}

func (dr *depositRepository) GetNewDepositByAccountIdUserId(accountId, userId string) (*Deposit, error) {
	var deposit Deposit
	queryParams := Deposit{
		AccountID:      accountId,
		UserID:         userId,
		ApprovalStatus: enums.DepositApprovalStatusNew,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Where(&queryParams).First(&deposit).Error; err != nil {
		return nil, err
	}

	return &deposit, nil
}

func (dr *depositRepository) GetPendingAccountByAccountIdUserId(accountId, userId string) (*Deposit, error) {
	var deposit Deposit
	queryParams := Deposit{
		AccountID:      accountId,
		UserID:         userId,
		ApprovalStatus: enums.DepositApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Where(&queryParams).First(&deposit).Error; err != nil {
		return nil, err
	}

	return &deposit, nil
}

func (dr *depositRepository) GetDepositsByUserIdAccountId(userId, accountId string) (*[]Deposits, error) {
	var deposits []Deposits
	if err := dr.db.Table("deposits").
		Joins("JOIN users ON users.id = deposits.user_id").
		Select("deposits.id as depositid, deposits.account_id as accountid, deposits.user_id as userid, users.name as name, users.email, deposits.amount, deposits.approval_status as approval_status, deposits.created_at as transaction_date").
		Where("deposits.approval_status != ? AND deposits.is_active = ? AND deposits.user_id = ? AND deposits.account_id = ?", 0, true, userId, accountId).
		Scan(&deposits).Error; err != nil {

		return nil, err
	}

	return &deposits, nil
}

func (dr *depositRepository) GetDepositByIdUserId(userId, depositId string) (*Deposit, error) {
	var deposit Deposit
	queryParams := Deposit{
		ID:     depositId,
		UserID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Where(&queryParams).First(&deposit).Error; err != nil {
		return nil, err
	}

	return &deposit, nil
}

// func (dr *depositRepository) GetDepositDetail(userId, depositId string) (*DepositDetail, error) {
// 	var depositDetail DepositDetail
// 	if err := dr.db.Table("deposits").
// 		Joins("JOIN users ON users.id = deposits.user_id").
// 		Select("deposits.id as depositid, deposits.account_id as accountid, deposits.user_id as userid, users.name as name, users.email, deposits.amount, deposits.bank_name, deposits.bank_account_number, deposits.bank_beneficiary_name, deposits.deposit_photo as deposit_photo, deposits.approval_status as approval_status").
// 		Where("deposits.approval_status != ? AND deposits.is_active = ? AND deposits.id = ? AND deposits.user_id = ?", 0, true, depositId, userId).
// 		Scan(&depositDetail).Error; err != nil {

// 		return nil, err
// 	}

// 	return &depositDetail, nil
// }

func (dr *depositRepository) GetBackOfficePendingDeposits() (*[]BackOfficePendingDeposit, error) {
	var backOfficePendingDeposits []BackOfficePendingDeposit
	if err := dr.db.Table("deposits").
		Joins("JOIN users ON users.id = deposits.user_id").
		Select("deposits.id as depositid, deposits.account_id as accountid, deposits.user_id as userid, users.name as name, users.email, deposits.amount, deposits.approval_status as approval_status").
		Where("deposits.approval_status = ? AND deposits.is_active = ?", enums.DepositApprovalStatusPending, true).
		Scan(&backOfficePendingDeposits).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingDeposits, nil
}

func (dr *depositRepository) GetBackOfficePendingDepositDetail(depositId string) (*BackOfficePendingDepositDetail, error) {
	var backOfficePendingDepositDetail BackOfficePendingDepositDetail
	if err := dr.db.Table("deposits").
		Joins("JOIN users ON users.id = deposits.user_id").
		Select("deposits.id as depositid, deposits.account_id as accountid, deposits.user_id as userid, users.name as name, users.email, deposits.amount, deposits.bank_name, deposits.bank_account_number, deposits.bank_beneficiary_name, deposits.deposit_photo as deposit_photo, deposits.approval_status as approval_status").
		Where("deposits.approval_status = ? AND deposits.is_active = ? AND deposits.id = ?", enums.DepositApprovalStatusPending, true, depositId).
		Scan(&backOfficePendingDepositDetail).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingDepositDetail, nil
}

func (dr *depositRepository) GetPendingDepositsById(depositId string) (*Deposit, error) {
	var deposits Deposit
	queryParams := Deposit{
		ID:             depositId,
		ApprovalStatus: enums.DepositApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := dr.db.Find(&deposits, &queryParams).Error; err != nil {
		return nil, err
	}

	return &deposits, nil
}

func (dr *depositRepository) UpdateDepositApprovalStatus(deposit *Deposit) error {
	return dr.db.Select(
		"DepositType",
		"ApprovalStatus",
		"ApprovedBy",
		"ApprovedAt",
	).Updates(deposit).Error
}

func (dr *depositRepository) Update(deposit *Deposit) error {
	return dr.db.Updates(deposit).Error
}

func (dr *depositRepository) SaveDepositPhoto(deposit *Deposit) error {
	return dr.db.Save(deposit).Error
}
