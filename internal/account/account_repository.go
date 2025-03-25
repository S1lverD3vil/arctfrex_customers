package account

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (ar *accountRepository) Create(account *Account) error {
	return ar.db.Create(account).Error
}

func (ar *accountRepository) GetPendingAccountByUserId(userId string) (*Account, error) {
	var account Account
	queryParams := Account{
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

func (ar *accountRepository) GetPendingAccountsByUserdId(userId string) (*[]Account, error) {
	var accounts []Account
	queryParams := Account{
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

func (ar *accountRepository) GetPendingAccountsById(accountId string) (*Account, error) {
	var accounts Account
	queryParams := Account{
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

func (ar *accountRepository) GetAccountsByUserdId(userId string) (*[]Account, error) {
	var accounts []Account
	queryParams := Account{
		UserID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Where(&queryParams).Not(&Account{ApprovalStatus: enums.AccountApprovalStatusCancelled}).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (ar *accountRepository) GetAccountsById(accountId string) (*Account, error) {
	var account Account
	queryParams := Account{
		ID: accountId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Where(&queryParams).Not(&Account{ApprovalStatus: enums.AccountApprovalStatusCancelled}).Find(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (ar *accountRepository) GetAccountsByIdUserId(userId, accountId string) (*Account, error) {
	var account Account
	queryParams := Account{
		ID:     accountId,
		UserID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ar.db.Where(&queryParams).Not(&Account{ApprovalStatus: enums.AccountApprovalStatusCancelled}).Find(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (ar *accountRepository) GetPendingAccounts() (*[]Account, error) {
	var accounts []Account
	queryParams := Account{
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

func (ar *accountRepository) GetBackOfficePendingAccountUserData(userid string) (*AccountUserData, error) {
	var accountUserData AccountUserData
	if err := ar.db.Table("users").
		// Joins("JOIN users ON users.id = accounts.user_id").
		Select("users.id as userid, users.name as name, users.email, users.mobile_phone as mobile_phone").
		Where("users.id = ? AND users.is_active = ?", userid, true).
		Scan(&accountUserData).Error; err != nil {

		return nil, err
	}

	return &accountUserData, nil
}

func (ar *accountRepository) GetBackOfficePendingAccounts() (*[]BackOfficePendingAccount, error) {
	var backOfficePendingAccounts []BackOfficePendingAccount
	if err := ar.db.Table("accounts").
		Joins("JOIN users ON users.id = accounts.user_id").
		Select("accounts.id as accountid, accounts.user_id as userid, users.name as name, users.email, accounts.approval_status as approval_status").
		Where("accounts.approval_status = ? AND accounts.is_active = ?", enums.AccountApprovalStatusPending, true).
		Scan(&backOfficePendingAccounts).Error; err != nil {

		return nil, err
	}

	return &backOfficePendingAccounts, nil
}

func (ar *accountRepository) GetBackOfficeAllAccounts() (*[]BackOfficeAllAccount, error) {
	var backOfficeAllAccounts []BackOfficeAllAccount
	if err := ar.db.Table("accounts").
		Joins("JOIN users ON users.id = accounts.user_id").
		Select(`
			accounts.id as accountid,
			accounts.user_id as userid,
			users.name as name,
			users.email,
			accounts.approval_status as approval_status
		`).
		Where(`
			accounts.is_active = ?
			AND accounts.type = ?`,
			true,
			enums.AccountTypeReal,
		).
		Scan(&backOfficeAllAccounts).Error; err != nil {

		return nil, err
	}

	return &backOfficeAllAccounts, nil
}

func (ar *accountRepository) GetReportProfitLoss(startDate, endDate string) (*[]ReportProfitLossData, error) {
	var reportProfitLossData []ReportProfitLossData
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

func (ar *accountRepository) UpdateAccount(account *Account) error {
	return ar.db.Updates(account).Error
}

func (ar *accountRepository) UpdateAccountApprovalStatus(account *Account) error {
	return ar.db.Select(
		"ApprovalStatus",
		"ApprovedBy",
		"ApprovedAt",
		"MetaLoginId",
		"MetaLoginPassword",
	).Updates(account).Error
}

func (ar *accountRepository) UpdateRealAccountCallRecording(account *Account) error {
	return ar.db.Select("RealAccountCallRecording").Updates(account).Error
}
