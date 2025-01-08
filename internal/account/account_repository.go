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
