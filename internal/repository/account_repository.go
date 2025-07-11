package repository

import (
	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
)

type AccountRepository interface {
	Create(account *model.Account) error
	GetPendingAccountByUserId(userId string) (*model.Account, error)
	GetPendingAccountsByUserdId(userId string) (*[]model.Account, error)
	GetPendingAccountsById(accountId string) (*model.Account, error)
	GetAccountsByUserdId(userId string) (*[]model.Account, error)
	GetAccountsById(accountId string) (*model.Account, error)
	GetAccountsByIdUserId(userId, accountId string) (*model.Account, error)
	GetPendingAccounts() (*[]model.Account, error)
	GetBackOfficePendingAccountUserData(userid string) (*model.AccountUserData, error)
	GetBackOfficePendingAccounts(request model.BackOfficePendingAccountRequest) ([]model.BackOfficePendingAccount, error)
	GetBackOfficeAllAccounts(request model.BackOfficeAllAccountRequest) ([]model.BackOfficeAllAccount, error)
	GetBackOfficeAccountByFilterParams(request model.BackOfficeAccountByFilterParams) ([]model.BackOfficeAccountByFilterParamsResponse, error)
	GetReportProfitLoss(startDate, endDate string) (*[]model.ReportProfitLossData, error)
	UpdateAccount(account *model.Account) error
	UpdateAccountApprovalStatus(account *model.Account) error
	UpdateRealAccountCallRecording(account *model.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (ar *accountRepository) Create(account *model.Account) error {
	return ar.db.Create(account).Error
}

func (ar *accountRepository) GetPendingAccountByUserId(userId string) (*model.Account, error) {
	var account model.Account
	queryParams := model.Account{
		UserID:         userId,
		ApprovalStatus: enums.AccountApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Where(&queryParams).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (ar *accountRepository) GetPendingAccountsByUserdId(userId string) (*[]model.Account, error) {
	var accounts []model.Account
	queryParams := model.Account{
		UserID:         userId,
		ApprovalStatus: enums.AccountApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.First(&accounts, &queryParams).Error; err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (ar *accountRepository) GetPendingAccountsById(accountId string) (*model.Account, error) {
	var accounts model.Account
	queryParams := model.Account{
		ID:             accountId,
		ApprovalStatus: enums.AccountApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Find(&accounts, &queryParams).Error; err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (ar *accountRepository) GetAccountsByUserdId(userId string) (*[]model.Account, error) {
	var accounts []model.Account
	queryParams := model.Account{
		UserID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Where(&queryParams).Not(&model.Account{ApprovalStatus: enums.AccountApprovalStatusCancelled}).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (ar *accountRepository) GetAccountsById(accountId string) (*model.Account, error) {
	var account model.Account
	queryParams := model.Account{
		ID: accountId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Where(&queryParams).Not(&model.Account{ApprovalStatus: enums.AccountApprovalStatusCancelled}).Find(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (ar *accountRepository) GetAccountsByIdUserId(userId, accountId string) (*model.Account, error) {
	var account model.Account
	queryParams := model.Account{
		ID:     accountId,
		UserID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Where(&queryParams).Not(&model.Account{ApprovalStatus: enums.AccountApprovalStatusCancelled}).Find(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (ar *accountRepository) GetPendingAccounts() (*[]model.Account, error) {
	var accounts []model.Account
	queryParams := model.Account{
		ApprovalStatus: enums.AccountApprovalStatusPending,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Find(&accounts, &queryParams).Error; err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (ar *accountRepository) GetBackOfficePendingAccountUserData(userid string) (*model.AccountUserData, error) {
	var accountUserData model.AccountUserData
	if err := ar.db.Table("users").
		// Joins("JOIN users ON users.id = accounts.user_id").
		Select("users.id as userid, users.name as name, users.email, users.mobile_phone as mobile_phone").
		Where("users.id = ? AND users.is_active = ?", userid, true).
		Scan(&accountUserData).Error; err != nil {

		return nil, err
	}

	return &accountUserData, nil
}

func (ar *accountRepository) GetBackOfficePendingAccounts(request model.BackOfficePendingAccountRequest) ([]model.BackOfficePendingAccount, error) {
	var backOfficePendingAccounts []model.BackOfficePendingAccount

	query := ar.db.Table("accounts").
		Joins("JOIN users ON users.id = accounts.user_id").
		Select(`
			accounts.id AS account_id,
			accounts.user_id AS user_id,
			users.name AS name, users.email,
			accounts.approval_status AS approval_status,
			users.mobile_phone AS user_mobile_phone,
			users.fax_phone AS user_fax_phone,
			users.home_phone AS user_home_phone,
			accounts.no_aggreement,
			accounts.created_at
		`).
		Where("accounts.approval_status = ? AND accounts.is_active = ?", enums.AccountApprovalStatusPending, true)

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize
	if err := query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficePendingAccounts).Error; err != nil {
		return nil, err
	}

	return backOfficePendingAccounts, nil
}

func (ar *accountRepository) GetBackOfficeAllAccounts(request model.BackOfficeAllAccountRequest) ([]model.BackOfficeAllAccount, error) {
	var backOfficeAllAccounts []model.BackOfficeAllAccount

	query := ar.db.Table("accounts").
		Joins("JOIN users ON users.id = accounts.user_id").
		Select(`
			accounts.id AS account_id,
			accounts.user_id AS user_id,
			users.name AS name,
			users.email,
			accounts.approval_status AS approval_status,
			users.mobile_phone AS user_mobile_phone,
			users.fax_phone AS user_fax_phone,
			users.home_phone AS user_home_phone,
			accounts.no_aggreement,
			accounts.created_at
		`).
		Where(`
			accounts.is_active = ?
			AND accounts.type = ?`,
			true,
			enums.AccountTypeReal)

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize
	if err := query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficeAllAccounts).Error; err != nil {
		return nil, err
	}

	return backOfficeAllAccounts, nil
}

func (ar *accountRepository) GetBackOfficeAccountByFilterParams(request model.BackOfficeAccountByFilterParams) ([]model.BackOfficeAccountByFilterParamsResponse, error) {
	var backOfficeAccounts []model.BackOfficeAccountByFilterParamsResponse

	query := ar.db.Table("accounts").
		Joins("JOIN users ON users.id = accounts.user_id").
		Select(`
			accounts.id AS account_id,
			accounts.user_id AS user_id,
			users.name AS name,
			users.email,
			accounts.approval_status AS approval_status,
			users.mobile_phone AS user_mobile_phone,
			users.fax_phone AS user_fax_phone,
			users.home_phone AS user_home_phone,
			accounts.no_aggreement,
			accounts.created_at
		`).
		Where(`
			accounts.is_active = ?
			AND accounts.type = ?
			AND accounts.approval_status = ?`,
			true,
			request.Type.EnumIndex(),
			request.ApprovalStatus.EnumIndex(),
		)

	offset := (request.Pagination.CurrentPage - 1) * request.Pagination.PageSize
	if err := query.Count(&request.Pagination.Paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Limit(request.Pagination.PageSize).
		Offset(offset).
		Scan(&backOfficeAccounts).Error; err != nil {
		return nil, err
	}

	return backOfficeAccounts, nil
}

func (ar *accountRepository) GetReportProfitLoss(startDate, endDate string) (*[]model.ReportProfitLossData, error) {
	var reportProfitLossData []model.ReportProfitLossData
	query := ar.db.Table("accounts").
		Joins("JOIN users ON users.id = accounts.user_id").
		Joins("JOIN user_addresses ua ON ua.id = users.id").
		Joins("JOIN user_finances uf ON uf.id = users.id").
		Joins("LEFT JOIN (SELECT account_id, user_id, SUM(amount) AS total_amount FROM deposits GROUP BY account_id, user_id) total_deposit ON total_deposit.account_id = accounts.id AND total_deposit.user_id = users.id").
		Joins("LEFT JOIN (SELECT account_id, user_id, SUM(amount) AS total_amount FROM withdrawals GROUP BY account_id, user_id) total_withdrawal ON total_withdrawal.account_id = accounts.id AND total_withdrawal.user_id = users.id").
		Select(`
			accounts.meta_login_id,
			users.name,
			ua.dom_city,
			uf.currency,
			uf.currency_rate,
			COALESCE(total_deposit.total_amount, 0) AS total_deposit_amount,
			COALESCE(total_withdrawal.total_amount, 0) AS total_withdrawal_amount,
			COALESCE(total_deposit.total_amount, 0) - COALESCE(total_withdrawal.total_amount, 0) AS prev_equity,
			COALESCE(total_deposit.total_amount, 0) - COALESCE(total_withdrawal.total_amount, 0) AS nmii,
			accounts.equity AS last_equity,
			accounts.equity - 0 AS gross_profit,
			accounts.equity - 0 AS gross_profit_usd,
			0 AS single_side_lot,
			0 AS commission,
			0 AS rebate,
			0 AS prev_bad_debt,
			0 AS last_bad_debt,
			accounts.equity - 0 AS net_profit,
			accounts.equity - 0 AS net_profit_usd,
			accounts.id AS accountid,
			users.id AS userid
		`).
		Where(`
			accounts.type = ?
			and users.mobile_phone = ?`,
			enums.AccountTypeReal,
			"812982951181",
		)

	if startDate != "" && endDate != "" {
		query = query.Where("accounts.created_at BETWEEN ? AND ?", startDate, endDate)

	}

	if err := query.Scan(&reportProfitLossData).Error; err != nil {
		return nil, err
	}

	return &reportProfitLossData, nil
}

func (ar *accountRepository) UpdateAccount(account *model.Account) error {
	return ar.db.Updates(account).Error
}

func (ar *accountRepository) UpdateAccountApprovalStatus(account *model.Account) error {
	return ar.db.Select(
		"ApprovalStatus",
		"ApprovedBy",
		"ApprovedAt",
		"MetaLoginId",
		"MetaLoginPassword",
	).Updates(account).Error
}

func (ar *accountRepository) UpdateRealAccountCallRecording(account *model.Account) error {
	return ar.db.Select("RealAccountCallRecording").Updates(account).Error
}
