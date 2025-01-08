package account

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"
	"time"
)

type Account struct {
	ID                       string                      `json:"accountid" gorm:"primary_key"`
	UserID                   string                      `json:"userid"`
	Type                     enums.AccountType           `json:"type"`
	Balance                  float64                     `json:"balance"`
	Equity                   float64                     `json:"equity"`
	Credit                   float64                     `json:"credit"`
	Margin                   float64                     `json:"margin"`
	MarginLevel              float64                     `json:"margin_level"`
	FreeMargin               float64                     `json:"free_margin"`
	WinTrades                float64                     `json:"win_trades"`
	TotalPL                  float64                     `json:"total_pl"`
	ApprovalStatus           enums.AccountApprovalStatus `json:"approval_status"`
	ApprovedBy               string                      `json:"approved_by"`
	ApprovedAt               time.Time                   `json:"approved_at"`
	IsDemo                   bool                        `json:"is_demo"`
	MetaLoginId              int64                       `json:"meta_login_id"`
	MetaLoginPassword        string                      `json:"meta_login_password"`
	RealAccountCallRecording string                      `json:"realaccount_callrecording"`

	base.BaseModel
}

type TopUpAccount struct {
	ID     string  `json:"accountid"`
	UserID string  `json:"userid"`
	Amount float64 `json:"amount"`
}

type BackOfficeAllAccount struct {
	Accountid      string                      `json:"accountid"`
	Userid         string                      `json:"userid"`
	Name           string                      `json:"name"`
	Email          string                      `json:"email"`
	ApprovalStatus enums.AccountApprovalStatus `json:"approval_status"`
}

type BackOfficePendingAccount struct {
	Accountid      string                      `json:"accountid"`
	Userid         string                      `json:"userid"`
	Name           string                      `json:"name"`
	Email          string                      `json:"email"`
	ApprovalStatus enums.AccountApprovalStatus `json:"approval_status"`
}

type BackOfficePendingApprovalRequest struct {
	Accountid string `json:"accountid"`
	Decision  string `json:"decision"`
	UserLogin string `json:"userlogin"`
}

type AccountApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}

type ClientAdd struct {
	Login    int64  `json:"Login"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
	Group    string `json:"Group"`
	Leverage int64  `json:"Leverage"`
	Rights   int64  `json:"Rights"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
}

type AccountUserData struct {
	Userid      string `json:"userid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobile_phone"`
}

type DemoAccountTopUp struct {
	Login  int64   `json:"Login"`
	Amount float64 `json:"Amount"`
	Result string  `json:"result"`
}

type AccountRepository interface {
	Create(account *Account) error
	GetPendingAccountByUserId(userId string) (*Account, error)
	GetPendingAccountsByUserdId(userId string) (*[]Account, error)
	GetPendingAccountsById(accountId string) (*Account, error)
	GetAccountsByUserdId(userId string) (*[]Account, error)
	GetAccountsById(accountId string) (*Account, error)
	GetAccountsByIdUserId(userId, accountId string) (*Account, error)
	GetPendingAccounts() (*[]Account, error)
	GetBackOfficePendingAccountUserData(userid string) (*AccountUserData, error)
	GetBackOfficePendingAccounts() (*[]BackOfficePendingAccount, error)
	GetBackOfficeAllAccounts() (*[]BackOfficeAllAccount, error)
	UpdateAccount(account *Account) error
	UpdateAccountApprovalStatus(account *Account) error
	UpdateRealAccountCallRecording(account *Account) error
}
