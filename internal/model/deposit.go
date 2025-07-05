package model

import (
	"errors"
	"slices"
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
	CreditType              enums.CreditType            `json:"credit_type"`
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
	IsInitialMargin     bool                        `json:"is_initial_margin"`
}

type BackOfficePendingApprovalRequest struct {
	Depositid   string            `json:"depositid"`
	DepositType enums.DepositType `json:"deposit_type"`
	Decision    string            `json:"decision"`
	UserLogin   string            `json:"userlogin"`
}

type BackOfficeUpdateCreditTypeRequest struct {
	Depositid           string `json:"deposit_id"`
	CreditTypeLocaleKey string `json:"credit_type_locale_key"`
}

func (b BackOfficeUpdateCreditTypeRequest) Validate() error {
	if b.Depositid == common.STRING_EMPTY {
		return errors.New("deposit ID is required")
	}

	if b.CreditTypeLocaleKey == common.STRING_EMPTY {
		return errors.New("credit type locale key is required")
	}

	if _, ok := enums.CreditTypeLocaleKeyToId[b.CreditTypeLocaleKey]; !ok {
		return errors.New("invalid credit type locale key")
	}

	return nil
}

type ApiResponse struct {
	base.ApiResponse
}

type ApiPaginatedResponse struct {
	base.ApiPaginatedResponse
}

type TradeDepositRequest struct {
	Login     int64   `json:"Login"`
	Amount    float64 `json:"Amount"`
	Result    string  `json:"result"`
	TradeType string  `json:"trade_type"`
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

type CreditBackOfficeDetailParam struct {
	Menutype  string
	DepositID string
}

func (c CreditBackOfficeDetailParam) Validate() error {
	if c.Menutype == common.STRING_EMPTY {
		return errors.New("menutype is required")
	}

	if c.DepositID == common.STRING_EMPTY {
		return errors.New("deposit ID is required")
	}

	if !slices.Contains([]string{common.Finance, common.Settlement}, c.Menutype) {
		return errors.New("invalid menutype")
	}

	return nil
}

type BackOfficeCreditDetail struct {
	DepositID           string                      `json:"deposit_id"`
	AccountID           string                      `json:"account_id"`
	UserID              string                      `json:"user_id"`
	Name                string                      `json:"name"`
	Email               string                      `json:"email"`
	Amount              float64                     `json:"amount"`
	AmountUsd           float64                     `json:"amount_usd"`
	BankName            string                      `json:"bank_name"`
	BankAccountNumber   string                      `json:"bank_account_number"`
	BankBeneficiaryName string                      `json:"bank_beneficiary_name"`
	DepositPhoto        string                      `json:"deposit_photo"`
	ApprovalStatus      enums.DepositApprovalStatus `json:"approval_status"`
	IsInitialMargin     bool                        `json:"is_initial_margin"`
}
