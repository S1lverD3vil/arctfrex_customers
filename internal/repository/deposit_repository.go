package repository

import (
	"fmt"

	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
)

type DepositRepository interface {
	Create(deposit *model.Deposit) error
	GetNewDepositByAccountIdUserId(accountId, userId string) (*model.Deposit, error)
	GetPendingAccountByAccountIdUserId(accountId, userId string) (*model.Deposit, error)
	GetDepositsByUserIdAccountId(userId, accountId string) (*[]model.Deposits, error)
	GetDepositsByUserIDAccountIDForIsInitialMargin(userId, accountId string) (*model.InitialMargin, error)
	GetDepositByIdUserId(userId, depositId string) (*model.Deposit, error)
	GetBackOfficePendingDeposits() (*[]model.BackOfficePendingDeposit, error)
	GetBackOfficePendingDepositDetail(depositId string) (*model.BackOfficePendingDepositDetail, error)
	GetPendingDepositsById(depositId string) (*model.Deposit, error)
	UpdateDepositApprovalStatus(deposit *model.Deposit) error
	Update(deposit *model.Deposit) error
	SaveDepositPhoto(deposit *model.Deposit) error
	GetBackOfficePendingDepositSPA(request model.DepositBackOfficeParam) ([]model.BackOfficePendingDeposit, error)
	GetBackOfficePendingDepositMulti(request model.DepositBackOfficeParam) ([]model.BackOfficePendingDeposit, error)
	GetBackOfficeCreditSPA(request model.CreditBackOfficeParam) ([]model.BackOfficeCreditInOut, error)
	GetBackOfficeCreditMulti(request model.CreditBackOfficeParam) ([]model.BackOfficeCreditInOut, error)
	GetBackOfficeCreditDetailByDepositID(depositId string, creditType enums.CreditType) (*model.BackOfficeCreditDetail, error)
}

type depositRepository struct {
	db *gorm.DB
}

func NewDepositRepository(db *gorm.DB) DepositRepository {
	return &depositRepository{db: db}
}

func (dr *depositRepository) Create(deposit *model.Deposit) error {
	return dr.db.Create(deposit).Error
}

func (dr *depositRepository) GetNewDepositByAccountIdUserId(accountId, userId string) (*model.Deposit, error) {
	var deposit model.Deposit
	queryParams := model.Deposit{
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

func (dr *depositRepository) GetPendingAccountByAccountIdUserId(accountId, userId string) (*model.Deposit, error) {
	var deposit model.Deposit
	queryParams := model.Deposit{
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

func (dr *depositRepository) GetDepositsByUserIdAccountId(userId, accountId string) (*[]model.Deposits, error) {
	var deposits []model.Deposits
	if err := dr.db.Table("deposits").
		Joins("JOIN users ON users.id = deposits.user_id").
		Select("deposits.id as depositid, deposits.account_id as accountid, deposits.user_id as userid, users.name as name, users.email, deposits.amount, deposits.approval_status as approval_status, deposits.created_at as transaction_date").
		Where("deposits.approval_status != ? AND deposits.is_active = ? AND deposits.user_id = ? AND deposits.account_id = ?", 0, true, userId, accountId).
		Scan(&deposits).Error; err != nil {

		return nil, err
	}

	return &deposits, nil
}

func (dr *depositRepository) GetDepositsByUserIDAccountIDForIsInitialMargin(userId, accountId string) (*model.InitialMargin, error) {
	var initialMargin model.InitialMargin

	if err := dr.db.Table("accounts").
		Select("accounts.type AS account_type, COUNT(deposits.id) AS total").
		Joins(`LEFT JOIN deposits 
			ON deposits.account_id = accounts.id 
			AND deposits.user_id = accounts.user_id 
			AND deposits.is_active = ? 
			AND deposits.approval_status = ?`,
			true, enums.DepositApprovalStatusApproved).
		Where("accounts.type = ? AND accounts.id = ? AND accounts.user_id = ? AND accounts.is_active = ?",
			enums.AccountTypeReal, accountId, userId, true).
		Group("accounts.type").
		Scan(&initialMargin).Error; err != nil {

		return &initialMargin, err
	}

	return &initialMargin, nil
}

func (dr *depositRepository) GetDepositByIdUserId(userId, depositId string) (*model.Deposit, error) {
	var deposit model.Deposit
	queryParams := model.Deposit{
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

func (dr *depositRepository) GetBackOfficePendingDeposits() (*[]model.BackOfficePendingDeposit, error) {
	var backOfficePendingDeposits []model.BackOfficePendingDeposit
	if err := dr.db.Table("deposits").
		Joins("JOIN users ON users.id = deposits.user_id").
		Select("deposits.id as depositid, deposits.account_id as accountid, deposits.user_id as userid, users.name as name, users.email, deposits.amount, deposits.approval_status as approval_status").
		Where("deposits.approval_status = ? AND deposits.is_active = ?", enums.DepositApprovalStatusPending, true).
		Scan(&backOfficePendingDeposits).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingDeposits, nil
}

func (dr *depositRepository) GetBackOfficePendingDepositSPA(request model.DepositBackOfficeParam) (backOfficePendingDeposits []model.BackOfficePendingDeposit, err error) {
	baseSelect := `
		deposits.id AS deposit_id,
		deposits.account_id,
		deposits.user_id, 
		users.name AS name,
		users.email,
		users.meta_login_id, 
		deposits.amount, 
		deposits.approval_status,
		deposits.created_at,
		deposits.bank_name,
		deposits.deposit_type
	`

	if request.Menutype == common.Settlement {
		baseSelect += `,
		bu.name AS finance_by`
	}

	query := dr.db.Table("deposits").
		Select(baseSelect).
		Joins("JOIN users ON users.id = deposits.user_id").
		Where("deposits.approval_status = ? AND deposits.is_active = ?", enums.DepositApprovalStatusPending, true)

	switch request.Menutype {
	case common.Finance:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id").
			Where("wa1.level=1 AND wa1.is_active=? AND wa1.status=?", true, enums.AccountApprovalStatusPending).
			Where("deposits.credit_type IS NULL OR deposits.credit_type = ?", enums.TypeCreditDefault)
	case common.Settlement:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id AND wa1.level = 1 AND (wa1.status = ? OR deposits.credit_type = ?) AND wa1.is_active = ?", enums.AccountApprovalStatusApproved, enums.TypeCreditIn, true).
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = deposits.id").
			Where("wa2.level=2 AND wa2.is_active=? AND wa2.status=?", true, enums.AccountApprovalStatusPending)
	default:
		return backOfficePendingDeposits, fmt.Errorf("invalid menu type: %s", request.Menutype)
	}

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize

	if err = query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err = query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficePendingDeposits).Error; err != nil {
		return nil, err
	}

	return backOfficePendingDeposits, nil
}

func (dr *depositRepository) GetBackOfficePendingDepositMulti(request model.DepositBackOfficeParam) (backOfficePendingDeposits []model.BackOfficePendingDeposit, err error) {
	baseSelect := `
		deposits.id AS deposit_id,
		deposits.account_id,
		deposits.user_id, 
		users.name AS name,
		users.email,
		users.meta_login_id, 
		deposits.amount, 
		deposits.approval_status,
		deposits.created_at,
		deposits.bank_name,
		deposits.deposit_type
	`

	if request.Menutype == common.Settlement {
		baseSelect += `,
		bu.name AS finance_by`
	}

	query := dr.db.Table("deposits").
		Select(baseSelect).
		Joins("JOIN users ON users.id = deposits.user_id").
		Where("deposits.approval_status = ? AND deposits.is_active = ?", enums.DepositApprovalStatusPending, true)

	switch request.Menutype {
	case common.Finance:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id").
			Where("wa1.level=1 AND wa1.is_active=? AND wa1.status=?", true, enums.AccountApprovalStatusPending).
			Where("deposits.credit_type IS NULL OR deposits.credit_type = ?", enums.TypeCreditDefault)
	case common.Settlement:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id AND wa1.level = 1 AND (wa1.status = ? OR deposits.credit_type = ?) AND wa1.is_active = ?", enums.AccountApprovalStatusApproved, enums.TypeCreditIn, true).
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = deposits.id").
			Where("wa2.level=2 AND wa2.is_active=? AND wa2.status=?", true, enums.AccountApprovalStatusPending)
	default:
		return backOfficePendingDeposits, fmt.Errorf("invalid menu type: %s", request.Menutype)
	}

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize

	if err = query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err = query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficePendingDeposits).Error; err != nil {
		return nil, err
	}

	return backOfficePendingDeposits, nil
}

func (dr *depositRepository) GetBackOfficePendingDepositDetail(depositId string) (*model.BackOfficePendingDepositDetail, error) {
	var backOfficePendingDepositDetail model.BackOfficePendingDepositDetail
	if err := dr.db.Table("deposits").
		Joins("JOIN users ON users.id = deposits.user_id").
		Select("deposits.id as depositid, deposits.account_id as accountid, deposits.user_id as userid, users.name as name, users.email, deposits.amount,deposits.amount_usd, deposits.bank_name, deposits.bank_account_number, deposits.bank_beneficiary_name, deposits.deposit_photo as deposit_photo, deposits.approval_status as approval_status").
		Where("deposits.approval_status = ? AND deposits.is_active = ? AND deposits.id = ?", enums.DepositApprovalStatusPending, true, depositId).
		Scan(&backOfficePendingDepositDetail).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingDepositDetail, nil
}

func (dr *depositRepository) GetPendingDepositsById(depositId string) (*model.Deposit, error) {
	var deposits model.Deposit
	queryParams := model.Deposit{
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

func (dr *depositRepository) UpdateDepositApprovalStatus(deposit *model.Deposit) error {
	return dr.db.Select(
		"DepositType",
		"ApprovalStatus",
		"ApprovedBy",
		"ApprovedAt",
		"CreditType",
	).Updates(deposit).Error
}

func (dr *depositRepository) Update(deposit *model.Deposit) error {
	return dr.db.Updates(deposit).Error
}

func (dr *depositRepository) SaveDepositPhoto(deposit *model.Deposit) error {
	return dr.db.Save(deposit).Error
}

func (dr *depositRepository) GetBackOfficeCreditSPA(request model.CreditBackOfficeParam) (backOfficeCredits []model.BackOfficeCreditInOut, err error) {
	baseSelect := `
		deposits.id AS deposit_id,
		deposits.account_id,
		deposits.user_id, 
		users.name AS name,
		users.email,
		users.meta_login_id, 
		deposits.amount, 
		deposits.approval_status,
		deposits.created_at,
		deposits.bank_name,
		deposits.deposit_type,
		bu.name AS finance_by
	`

	if request.Menutype == common.CreditIn {
		baseSelect += `,
		bu2.name AS dealing_by`
	}

	query := dr.db.Table("deposits").
		Select(baseSelect).
		Joins("JOIN users ON users.id = deposits.user_id").
		Where("deposits.is_active = ?", true).
		Where("deposits.approval_status != ?", enums.DepositApprovalStatusRejected)

	switch request.Menutype {
	case common.CreditIn:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id").
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = deposits.id").
			Joins("LEFT JOIN backoffice_users AS bu2 ON wa2.approved_by=bu2.id and bu2.is_active=?", true).
			Where("wa1.level=1 AND wa1.is_active=?", true).
			Where("wa2.level=2 AND wa2.is_active=?", true).
			Where("deposits.credit_type = ?", enums.TypeCreditIn)
	case common.CreditOut:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id AND wa1.level = 1 AND wa1.is_active = ?", true).
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = deposits.id").
			Where("wa2.level=2 AND wa2.is_active=? ", true).
			Where("deposits.credit_type = ?", enums.TypeCreditOut)
	default:
		return backOfficeCredits, fmt.Errorf("invalid menu type: %s", request.Menutype)
	}

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize

	if err = query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err = query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficeCredits).Error; err != nil {
		return nil, err
	}

	return backOfficeCredits, nil
}

func (dr *depositRepository) GetBackOfficeCreditMulti(request model.CreditBackOfficeParam) (backOfficeCredits []model.BackOfficeCreditInOut, err error) {
	baseSelect := `
		deposits.id AS deposit_id,
		deposits.account_id,
		deposits.user_id, 
		users.name AS name,
		users.email,
		users.meta_login_id, 
		deposits.amount, 
		deposits.approval_status,
		deposits.created_at,
		deposits.bank_name,
		deposits.deposit_type,
		bu.name AS finance_by
	`

	if request.Menutype == common.CreditIn {
		baseSelect += `,
		bu2.name AS dealing_by`
	}

	query := dr.db.Table("deposits").
		Select(baseSelect).
		Joins("JOIN users ON users.id = deposits.user_id").
		Where("deposits.is_active = ?", true).
		Where("deposits.approval_status != ?", enums.DepositApprovalStatusRejected)

	switch request.Menutype {
	case common.CreditIn:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id").
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = deposits.id").
			Joins("LEFT JOIN backoffice_users AS bu2 ON wa2.approved_by=bu2.id and bu2.is_active=?", true).
			Where("wa1.level=1 AND wa1.is_active=?", true).
			Where("wa2.level=2 AND wa2.is_active=?", true).
			Where("deposits.credit_type = ?", enums.TypeCreditIn)
	case common.CreditOut:
		query = query.
			Joins("JOIN workflow_approvers AS wa1 ON wa1.document_id = deposits.id AND wa1.level = 1 AND wa1.is_active = ?", true).
			Joins("LEFT JOIN backoffice_users AS bu ON wa1.approved_by=bu.id and bu.is_active=?", true).
			Joins("JOIN workflow_approvers AS wa2 ON wa2.document_id = deposits.id").
			Where("wa2.level=2 AND wa2.is_active=? ", true).
			Where("deposits.credit_type = ?", enums.TypeCreditOut)
	default:
		return backOfficeCredits, fmt.Errorf("invalid menu type: %s", request.Menutype)
	}

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize

	if err = query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err = query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficeCredits).Error; err != nil {
		return nil, err
	}

	return backOfficeCredits, nil
}

func (dr *depositRepository) GetBackOfficeCreditDetailByDepositID(depositId string, creditType enums.CreditType) (*model.BackOfficeCreditDetail, error) {
	var backOfficeCreditDetail model.BackOfficeCreditDetail
	if err := dr.db.Table("deposits").
		Joins("JOIN users ON users.id = deposits.user_id").
		Select("deposits.id as deposit_id, deposits.account_id , deposits.user_id , users.name as name, users.email, deposits.amount, deposits.amount_usd, deposits.bank_name, deposits.bank_account_number, deposits.bank_beneficiary_name, deposits.deposit_photo as deposit_photo, deposits.approval_status as approval_status").
		Where("deposits.is_active = ? AND deposits.id = ? AND deposits.credit_type = ?", true, depositId, creditType).
		Scan(&backOfficeCreditDetail).Error; err != nil {

		return nil, err
	}

	return &backOfficeCreditDetail, nil
}
