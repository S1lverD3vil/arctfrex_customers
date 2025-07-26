package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
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
	NoAggreement             string                      `json:"no_aggreement"`
	SurveyResult             SurveyResult                `gorm:"type:jsonb" json:"survey_result"` // Column for the array of obj

	base.BaseModel
}

type TopUpAccount struct {
	ID     string  `json:"accountid"`
	UserID string  `json:"userid"`
	Amount float64 `json:"amount"`
}

type BackOfficeAllAccount struct {
	AccountID       string                      `json:"account_id"`
	UserID          string                      `json:"user_id"`
	Name            string                      `json:"name"`
	Email           string                      `json:"email"`
	ApprovalStatus  enums.AccountApprovalStatus `json:"approval_status"`
	NoAggreement    string                      `json:"no_aggreement"`
	UserMobilePhone string                      `json:"user_mobile_phone"`
	UserFaxPhone    string                      `json:"user_fax_phone"`
	UserHomePhone   string                      `json:"user_home_phone"`
	CreatedAt       time.Time                   `json:"created_at"`
}

type BackOfficeAllAccountRequest struct {
	Pagination *common.TableListParams
}

type BackOfficeAllAccountResponse struct {
	Data       []BackOfficeAllAccount
	Pagination *common.TableListParams
}

type BackOfficePendingAccount struct {
	AccountID       string                      `json:"account_id"`
	UserID          string                      `json:"user_id"`
	Name            string                      `json:"name"`
	Email           string                      `json:"email"`
	ApprovalStatus  enums.AccountApprovalStatus `json:"approval_status"`
	NoAggreement    string                      `json:"no_aggreement"`
	UserMobilePhone string                      `json:"user_mobile_phone"`
	UserFaxPhone    string                      `json:"user_fax_phone"`
	UserHomePhone   string                      `json:"user_home_phone"`
	CreatedAt       time.Time                   `json:"created_at"`
}

type BackOfficePendingAccountRequest struct {
	Pagination *common.TableListParams
}

type BackOfficeAccountByMenuTypeRequest struct {
	MenuType       string                      `json:"menu_type" form:"menu_type"`
	Type           enums.AccountType           `json:"-"`
	ApprovalStatus enums.AccountApprovalStatus `json:"-"`
	Pagination     *common.TableListParams     `json:"pagination" form:"pagination"`
}

func (b *BackOfficeAccountByMenuTypeRequest) Normalize(accountType string, approvalStatus string) {
	// Normalize AccountType
	if val, ok := enums.AccountTypeMap[strings.ToLower(accountType)]; ok {
		b.Type = val
	} else {
		b.Type = enums.AccountTypeDemo
	}

	// Normalize ApprovalStatus
	if val, ok := enums.ApprovalStatusMap[strings.ToLower(approvalStatus)]; ok {
		b.ApprovalStatus = val
	} else {
		b.ApprovalStatus = enums.AccountApprovalStatusApproved
	}
}

type BackOfficeAccountByMenuTypeResponse struct {
	Data       []BackOfficeAccountByFilterParamsResponse
	Pagination *common.TableListParams
}

type BackOfficeAccountByFilterParams struct {
	Type           enums.AccountType           `json:"type"`
	ApprovalStatus enums.AccountApprovalStatus `json:"approval_status"`
	Pagination     *common.TableListParams
}

type BackOfficeAccountByFilterParamsResponse struct {
	AccountID       string                      `json:"account_id"`
	UserID          string                      `json:"user_id"`
	Name            string                      `json:"name"`
	Email           string                      `json:"email"`
	ApprovalStatus  enums.AccountApprovalStatus `json:"approval_status"`
	NoAggreement    string                      `json:"no_aggreement"`
	UserMobilePhone string                      `json:"user_mobile_phone"`
	UserFaxPhone    string                      `json:"user_fax_phone"`
	UserHomePhone   string                      `json:"user_home_phone"`
	CreatedAt       time.Time                   `json:"created_at"`
}

type BackOfficePendingAccountResponse struct {
	Data       []BackOfficePendingAccount
	Pagination *common.TableListParams
}

type BackOfficePendingAccountApprovalRequest struct {
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

type ReportProfitLossData struct {
	MetaLoginID           int64   `json:"meta_login_id"`
	Name                  string  `json:"name"`
	DomCity               string  `json:"dom_city"`
	Currency              string  `json:"currency"`
	CurrencyRate          float64 `json:"currency_rate"`
	TotalDepositAmount    float64 `json:"total_deposit_amount"`
	TotalWithdrawalAmount float64 `json:"total_withdrawal_amount"`
	PrevEquity            float64 `json:"prev_equity"`
	Nmii                  float64 `json:"nmii"`
	LastEquity            float64 `json:"last_equity"`
	GrossProfit           float64 `json:"gross_profit"`
	GrossProfitUSD        float64 `json:"gross_profit_usd"`
	SingleSideLot         float64 `json:"single_side_lot"`
	Commission            float64 `json:"commission"`
	Rebate                float64 `json:"rebate"`
	PrevBadDebt           float64 `json:"prev_bad_debt"`
	LastBadDebt           float64 `json:"last_bad_debt"`
	NetProfit             float64 `json:"net_profit"`
	NetProfitUSD          float64 `json:"net_profit_usd"`
	AccountID             int64   `json:"accountid"`
	UserID                int64   `json:"userid"`
}

type SurveyQuestionData struct {
	No     int  `json:"no"`
	Answer bool `json:"answer"`
}

type SurveyResult struct {
	SurveyChecklist []SurveyQuestionData `json:"survey_checklist"`
	PpatkChecklist  []SurveyQuestionData `json:"ppatk_checklist"`
}

type SurveyAnswerRequest struct {
	UserID          string               `json:"user_id"`
	AccountID       string               `json:"account_id"`
	SurveyChecklist []SurveyQuestionData `json:"survey_checklist"`
	PpatkChecklist  []SurveyQuestionData `json:"ppatk_checklist"`
}

type SurveyAnswerResponse struct {
	SurveyChecklist []SurveyQuestionData `json:"survey_checklist"`
	PpatkChecklist  []SurveyQuestionData `json:"ppatk_checklist"`
}

// Save to DB (jsonb)

// Marshal to JSONB
func (b SurveyResult) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Unmarshal from JSONB
func (b *SurveyResult) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, b)
}

type SurveyData struct {
	PhoneNumber            string       `json:"phone_number"`
	NoIdentity             string       `json:"no_identity"`
	Address                string       `json:"address"`
	DomProvince            string       `json:"dom_province"`
	Profession             string       `json:"profession"`
	MotherName             string       `json:"mother_name"`
	Email                  string       `json:"email"`
	Name                   string       `json:"name"`
	AccountType            string       `json:"account_type"`
	CurrencyRate           float64      `json:"currency_rate"`
	ProductServicePlatform string       `json:"product_service_platform"`
	AccountID              int          `json:"account_id"`
	BankList               BankList     `gorm:"type:jsonb" json:"bank_list"`
	SurveyResult           SurveyResult `gorm:"type:jsonb" json:"survey_result"`
}

type SurveyDataTemplate struct {
	PhoneNumber            string       `json:"phone_number"`
	NoIdentity             string       `json:"no_identity"`
	Address                string       `json:"address"`
	DomProvince            string       `json:"dom_province"`
	Profession             string       `json:"profession"`
	MotherName             string       `json:"mother_name"`
	Email                  string       `json:"email"`
	Name                   string       `json:"name"`
	AccountType            string       `json:"account_type"`
	CurrencyRate           float64      `json:"currency_rate"`
	ProductServicePlatform string       `json:"product_service_platform"`
	AccountID              int          `json:"account_id"`
	Banks                  BankList     `json:"banks"`
	SurveyResult           SurveyResult `json:"survey_result"`
}

type Banklist struct {
	Index               int    `json:"index"`
	BankName            string `json:"bank_name"`
	BankAccountNumber   string `json:"bank_account_number"`
	BankBeneficiaryName string `json:"bank_beneficiary_name"`
}
