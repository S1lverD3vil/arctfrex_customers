package model

import (
	"time"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
)

type Withdrawal struct {
	ID                      string                         `json:"withdrawalid" gorm:"primary_key"`
	Amount                  float64                        `json:"amount"`
	AmountUsd               float64                        `json:"amount_usd"`
	AccountID               string                         `json:"accountid"`
	UserID                  string                         `json:"userid"`
	PaymentMethod           string                         `json:"payment_method"`
	BankCode                string                         `json:"bank_code"`
	BankName                string                         `json:"bank_name"`
	BankAccountNumber       string                         `json:"bank_account_number"`
	BankBeneficiaryName     string                         `json:"bank_beneficiary_name"`
	SenderBankAccountNumber string                         `json:"sender_bank_account_number"`
	SenderBankAccountName   string                         `json:"sender_bank_account_name"`
	ApprovalStatus          enums.WithdrawalApprovalStatus `json:"approval_status"`
	ApprovedBy              string                         `json:"approved_by"`
	ApprovedAt              time.Time                      `json:"approved_at"`

	base.BaseModel
}

type Withdrawals struct {
	Withdrawalid    string                         `json:"withdrawalid"`
	Accountid       string                         `json:"accountid"`
	Userid          string                         `json:"userid"`
	Name            string                         `json:"name"`
	Email           string                         `json:"email"`
	Amount          float64                        `json:"amount"`
	AmountUsd       float64                        `json:"amount_usd"`
	ApprovalStatus  enums.WithdrawalApprovalStatus `json:"approval_status"`
	TransactionDate time.Time                      `json:"transaction_date"`
}

type BackOfficePendingWithdrawal struct {
	WithdrawalID   string                         `json:"withdrawal_id"`
	AccountID      string                         `json:"account_id"`
	UserID         string                         `json:"user_id"`
	Name           string                         `json:"name"`
	Email          string                         `json:"email"`
	Amount         float64                        `json:"amount"`
	AmountUSD      float64                        `json:"amount_usd"`
	ApprovalStatus enums.WithdrawalApprovalStatus `json:"approval_status"`
	CreatedAt      time.Time                      `json:"created_at"`
	BankName       string                         `json:"bank_name"`
	MetaLoginID    int64                          `json:"meta_login_id"`
	SettlementBy   string                         `json:"settlement_by"`
}

type BackOfficePendingWithdrawalDetail struct {
	Withdrawalid        string                         `json:"withdrawalid"`
	Accountid           string                         `json:"accountid"`
	Userid              string                         `json:"userid"`
	Name                string                         `json:"name"`
	Email               string                         `json:"email"`
	Amount              float64                        `json:"amount"`
	AmountUsd           float64                        `json:"amount_usd"`
	BankName            string                         `json:"bank_name"`
	BankAccountNumber   string                         `json:"bank_account_number"`
	BankBeneficiaryName string                         `json:"bank_beneficiary_name"`
	ApprovalStatus      enums.WithdrawalApprovalStatus `json:"approval_status"`
}

type BackOfficePendingWithdrawalApprovalRequest struct {
	Withdrawalid string `json:"withdrawalid"`
	Decision     string `json:"decision"`
	UserLogin    string `json:"userlogin"`
}

type WithdrawalApiResponse struct {
	base.ApiResponse
}

type TradeWithdrawal struct {
	Login  int64   `json:"Login"`
	Amount float64 `json:"Amount"`
	Result string  `json:"result"`
}

type WithdrawalBackOfficeParam struct {
	Menutype   string
	Pagination *common.TableListParams
}

type BackOfficePendingWithdrawalPagination struct {
	Data       []BackOfficePendingWithdrawal
	Pagination *common.TableListParams
}
