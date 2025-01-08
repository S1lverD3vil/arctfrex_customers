package deposit

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"
	"time"
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
	Depositid      string                      `json:"depositid"`
	Accountid      string                      `json:"accountid"`
	Userid         string                      `json:"userid"`
	Name           string                      `json:"name"`
	Email          string                      `json:"email"`
	Amount         float64                     `json:"amount"`
	AmountUsd      float64                     `json:"amount_usd"`
	ApprovalStatus enums.DepositApprovalStatus `json:"approval_status"`
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

type DepositApiResponse struct {
	base.ApiResponse
}

type TradeDeposit struct {
	Login  int64   `json:"Login"`
	Amount float64 `json:"Amount"`
	Result string  `json:"result"`
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
}
