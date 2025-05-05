package deposit

import (
	"time"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
)

type Deposit struct {
	ID                      string                      `json:"depositid" gorm:"primary_key"`
	Amount                  float64                     `json:"amount"`
	AmountUsd               float64                     `json:"amount_usd"`
	AccountID               string                      `json:"accountid"`
	UserID                  string                      `json:"userid"`
	PaymentMethod           string                      `json:"payment_method"`
	BankCode                string                      `json:"bank_code"`
	BankName                string                      `json:"bank_name"`
	BankAccountNumber       string                      `json:"bank_account_number"`
	BankBeneficiaryName     string                      `json:"bank_beneficiary_name"`
	SenderBankAccountNumber string                      `json:"sender_bank_account_number"`
	SenderBankAccountName   string                      `json:"sender_bank_account_name"`
	DepositPhoto            string                      `json:"deposit_photo"`
	DepositType             enums.DepositType           `json:"deposit_type"`
	ApprovalStatus          enums.DepositApprovalStatus `json:"approval_status"`
	ApprovedBy              string                      `json:"approved_by"`
	ApprovedAt              time.Time                   `json:"approved_at"`

	base.BaseModel
}

type Deposits struct {
	Depositid       string                      `json:"depositid"`
	Accountid       string                      `json:"accountid"`
	Userid          string                      `json:"userid"`
	Name            string                      `json:"name"`
	Email           string                      `json:"email"`
	Amount          float64                     `json:"amount"`
	AmountUsd       float64                     `json:"amount_usd"`
	ApprovalStatus  enums.DepositApprovalStatus `json:"approval_status"`
	TransactionDate time.Time                   `json:"transaction_date"`
}

type DepositDetail struct {
	Depositid           string                      `json:"depositid"`
	Accountid           string                      `json:"accountid"`
	Userid              string                      `json:"userid"`
	Name                string                      `json:"name"`
	Email               string                      `json:"email"`
	Amount              float64                     `json:"amount"`
	AmountUsd           float64                     `json:"amount_usd"`
	BankName            string                      `json:"bank_name"`
	BankAccountNumber   string                      `json:"bank_account_number"`
	BankBeneficiaryName string                      `json:"bank_beneficiary_name"`
	DepositPhoto        string                      `json:"deposit_photo"`
	ApprovalStatus      enums.DepositApprovalStatus `json:"approval_status"`
}

type BackOfficePendingDeposit struct {
	DepositID      string                      `json:"deposit_id"`
	AccountID      string                      `json:"account_id"`
	UserID         string                      `json:"user_id"`
	Name           string                      `json:"name"`
	Email          string                      `json:"email"`
	Amount         float64                     `json:"amount"`
	AmountUSD      float64                     `json:"amount_usd"`
	ApprovalStatus enums.DepositApprovalStatus `json:"approval_status"`
	CreatedAt      time.Time                   `json:"created_at"`
	BankName       string                      `json:"bank_name"`
	MetaLoginID    int64                       `json:"meta_login_id"`
	DepositType    enums.DepositType           `json:"deposit_type"`
	FinanceBy      string                      `json:"finance_by"`
}

type BackOfficeCreditInOut struct {
	DepositID      string                      `json:"deposit_id"`
	AccountID      string                      `json:"account_id"`
	UserID         string                      `json:"user_id"`
	Name           string                      `json:"name"`
	Email          string                      `json:"email"`
	Amount         float64                     `json:"amount"`
	AmountUSD      float64                     `json:"amount_usd"`
	ApprovalStatus enums.DepositApprovalStatus `json:"approval_status"`
	CreatedAt      time.Time                   `json:"created_at"`
	BankName       string                      `json:"bank_name"`
	MetaLoginID    int64                       `json:"meta_login_id"`
	DepositType    enums.DepositType           `json:"deposit_type"`
	FinanceBy      string                      `json:"finance_by"`
	DealingBy      string                      `json:"dealing_by"`
}

type BackOfficePendingDepositDetail struct {
	Depositid           string                      `json:"depositid"`
	Accountid           string                      `json:"accountid"`
	Userid              string                      `json:"userid"`
	Name                string                      `json:"name"`
	Email               string                      `json:"email"`
	Amount              float64                     `json:"amount"`
	AmountUsd           float64                     `json:"amount_usd"`
	BankName            string                      `json:"bank_name"`
	BankAccountNumber   string                      `json:"bank_account_number"`
	BankBeneficiaryName string                      `json:"bank_beneficiary_name"`
	DepositPhoto        string                      `json:"deposit_photo"`
	ApprovalStatus      enums.DepositApprovalStatus `json:"approval_status"`
}

type BackOfficePendingApprovalRequest struct {
	Depositid   string            `json:"depositid"`
	DepositType enums.DepositType `json:"deposit_type"`
	Decision    string            `json:"decision"`
	UserLogin   string            `json:"userlogin"`
}

type ApiResponse struct {
	base.ApiResponse
}

type ApiPaginatedResponse struct {
	base.ApiPaginatedResponse
}

type TradeDeposit struct {
	Login  int64   `json:"Login"`
	Amount float64 `json:"Amount"`
	Result string  `json:"result"`
}

type DepositBackOfficeParam struct {
	Menutype   string
	Pagination *common.TableListParams
}

type BackOfficePendingDepositPagination struct {
	Data       []BackOfficePendingDeposit
	Pagination *common.TableListParams
}

type CreditBackOfficeParam struct {
	Menutype   string
	Pagination *common.TableListParams
}

type BackOfficeCreditPagination struct {
	Data       []BackOfficeCreditInOut
	Pagination *common.TableListParams
}

type DepositRepository interface {
	Create(deposit *Deposit) error
	GetNewDepositByAccountIdUserId(accountId, userId string) (*Deposit, error)
	GetPendingAccountByAccountIdUserId(accountId, userId string) (*Deposit, error)
	GetDepositsByUserIdAccountId(userId, accountId string) (*[]Deposits, error)
	GetDepositByIdUserId(userId, depositId string) (*Deposit, error)
	GetBackOfficePendingDeposits() (*[]BackOfficePendingDeposit, error)
	GetBackOfficePendingDepositDetail(depositId string) (*BackOfficePendingDepositDetail, error)
	GetPendingDepositsById(depositId string) (*Deposit, error)
	UpdateDepositApprovalStatus(deposit *Deposit) error
	Update(deposit *Deposit) error
	SaveDepositPhoto(deposit *Deposit) error
	GetBackOfficePendingDepositSPA(request DepositBackOfficeParam) ([]BackOfficePendingDeposit, error)
	GetBackOfficePendingDepositMulti(request DepositBackOfficeParam) ([]BackOfficePendingDeposit, error)
	GetBackOfficeCreditSPA(request CreditBackOfficeParam) ([]BackOfficeCreditInOut, error)
	GetBackOfficeCreditMulti(request CreditBackOfficeParam) ([]BackOfficeCreditInOut, error)
}
