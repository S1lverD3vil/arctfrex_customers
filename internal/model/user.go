package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
)

type Users struct {
	ID                string         `json:"userid"`
	Name              string         `json:"customer_name"`
	Email             string         `gorm:"primary_key;uniqueIndex:idx_email_mobilephone" json:"email" binding:"required,email"`
	MobilePhone       string         `gorm:"primary_key;uniqueIndex:idx_email_mobilephone" json:"mobilephone" binding:"required"`
	HomePhone         string         `json:"home_phone"`
	FaxPhone          string         `json:"fax_phone"`
	Pin               string         `json:"pin"`
	Device            string         `json:"device"`
	DeviceId          string         `json:"device_id"`
	DeviceName        string         `json:"device_name"`
	DeviceImei        string         `json:"device_imei"`
	DeviceOs          string         `json:"device_os"`
	Latitude          string         `json:"latitude"`
	Longitude         string         `json:"longitude"`
	SessionId         string         `json:"session_id"`
	SessionExpiration time.Time      `json:"session_expiration"`
	Watchlist         pq.StringArray `gorm:"type:text[]" json:"watchlist"`
	MetaLoginId       int64          `json:"meta_login_id"`
	MetaLoginPassword string         `json:"meta_login_password"`
	ReferralCode      string         `json:"referral_code"`

	base.BaseModel
}

type UserProfile struct {
	ID                      string            `gorm:"primary_key" json:"userid"`
	Gender                  string            `json:"gender"`
	MaritalStatus           string            `json:"martial_status"`
	PlaceOfBirth            string            `json:"place_of_birth"`
	DateOfBirth             common.CustomDate `json:"date_of_birth"` // time_format:"2006-01-02"`
	Nationality             string            `json:"nationality"`
	MotherMaiden            string            `json:"mother_maiden"`
	KtpNumber               string            `json:"ktp_number"`
	NpwpNumber              string            `json:"npwp_number"`
	KtpPhoto                string            `json:"ktp_photo"`
	SelfiePhoto             string            `json:"selfie_photo"`
	NpwpPhoto               string            `json:"npwp_photo"`
	AdditionalDocumentPhoto string            `json:"additional_document_photo"`
	DeclarationVideo        string            `json:"declaration_video"`
	IdentityType            string            `json:"identity_type"` // e.g., KTP, Passport, etc.
	SpouseName              string            `json:"spouse_name"`
	DeclaredBankruptByCourt bool              `json:"declared_bankrupt_by_court"`
	FamilyAffiliation       bool              `json:"family_affiliation"`

	base.BaseModel
}

type UserProfileDetail struct {
	UserID                  string            `json:"user_id"`
	FullName                string            `json:"full_name"`
	MobilePhone             string            `json:"mobile_phone"`
	HomePhone               string            `json:"home_phone"`
	FaxPhone                string            `json:"fax_phone"`
	Gender                  string            `json:"gender"`
	MaritalStatus           string            `json:"martial_status"`
	PlaceOfBirth            string            `json:"place_of_birth"`
	DateOfBirth             common.CustomDate `json:"date_of_birth"` // time_format:"2006-01-02"`
	Nationality             string            `json:"nationality"`
	MotherMaiden            string            `json:"mother_maiden"`
	KTPNumber               string            `json:"ktp_number"`
	IdentityType            string            `json:"identity_type"` // e.g., KTP, Passport, etc.
	NPWPNumber              string            `json:"npwp_number"`
	KTPPhoto                string            `json:"ktp_photo"`
	SelfiePhoto             string            `json:"selfie_photo"`
	NPWPPhoto               string            `json:"npwp_photo"`
	AdditionalDocumentPhoto string            `json:"additional_document_photo"`
	DeclarationVideo        string            `json:"declaration_video"`
	DomPostalCode           string            `json:"dom_postal_code"`
}

type UserAddress struct {
	ID             string `gorm:"primary_key" json:"userid"`
	KtpCountry     string `json:"ktp_country"`
	KtpProvince    string `json:"ktp_province"`
	KtpCity        string `json:"ktp_city"`
	KtpDistrict    string `json:"ktp_district"`
	KtpSubDistrict string `json:"ktp_subdistrict"`
	KtpAddress     string `json:"ktp_address"`

	KtpSameDom bool `json:"ktp_same_dom"`

	DomCountry     string `json:"dom_country"`
	DomProvince    string `json:"dom_province"`
	DomCity        string `json:"dom_city"`
	DomDistrict    string `json:"dom_district"`
	DomSubDistrict string `json:"dom_subdistrict"`
	DomAddress     string `json:"dom_address"`
	DomPostalCode  string `json:"dom_postal_code"`

	ResidenceOwnership string `json:"residence_ownership"`

	base.BaseModel
}

type UserEmployment struct {
	ID                string `gorm:"primary_key" json:"userid"`
	CompanyName       string `json:"company_name"`
	CompanyAddress    string `json:"company_address"`
	CompanyCity       string `json:"company_city"`
	CompanyPhone      string `json:"company_phone"`
	CompanyPostalCode string `json:"company_postal_code"`
	WorkingSince      string `json:"working_since"`
	Profession        string `json:"profession"`
	WorkingField      string `json:"working_field"`
	PreviewJobTitle   string `json:"preview_job_title"`
	JobTitle          string `json:"job_title"`

	base.BaseModel
}

type UserEmploymentDetail struct {
	UserID            string `gorm:"primary_key" json:"user_id"`
	CompanyName       string `json:"company_name"`
	CompanyAddress    string `json:"company_address"`
	CompanyCity       string `json:"company_city"`
	CompanyPhone      string `json:"company_phone"`
	CompanyPostalCode string `json:"company_postal_code"`
	WorkingSince      string `json:"working_since"`
	Profession        string `json:"profession"`
	WorkingField      string `json:"working_field"`
	PreviewJobTitle   string `json:"preview_job_title"`
	JobTitle          string `json:"job_title"`
}

type UserFinance struct {
	ID                           string   `gorm:"primary_key" json:"userid"`
	SourceIncome                 string   `json:"source_income"`
	YearlyIncomeAmount           string   `json:"yearly_income_amount"`
	YearlyAdditionalIncomeAmount string   `json:"yearly_additional_income_amount"`
	EstimationWealthAmount       string   `json:"estimation_wealth_amount"`
	TaxableObjectSalesValue      string   `json:"taxable_object_sales_value"`
	Deposito                     string   `json:"deposito"`
	Currency                     string   `json:"currency"`
	InvestmentGoals              string   `json:"investment_goals"`
	InvestmentExperience         string   `json:"investment_experience"`
	BankList                     BankList `gorm:"type:jsonb" json:"bank_list"` // Column for the array of obj
	CurrencyRate                 float64  `json:"currency_rate"`
	ProductServicePlatform       string   `json:"product_service_platform"`
	AccountType                  string   `json:"account_type"` // e.g., Personal, Joint, etc.

	base.BaseModel
}

type Bank struct {
	BankName            string `json:"bank_name"`
	BankBranch          string `json:"bank_branch"`
	BankCity            string `json:"bank_city"`
	BankAccountNumber   string `json:"bank_account_number"`
	BankBeneficiaryName string `json:"bank_beneficiary_name"`
	BankAccountType     string `json:"bank_account_type"`
	BankPhone           string `json:"bank_phone"`
}

// Custom slice alias
type BankList []Bank

// Save to DB (jsonb)
func (b BankList) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Read from DB (jsonb)
func (b *BankList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to assert value to []byte")
	}
	return json.Unmarshal(bytes, b)
}

type UserFinanceDetail struct {
	UserID                       string   `json:"user_id"`
	SourceIncome                 string   `json:"source_income"`
	YearlyIncomeAmount           string   `json:"yearly_income_amount"`
	YearlyAdditionalIncomeAmount string   `json:"yearly_additional_income_amount"`
	EstimationWealthAmount       string   `json:"estimation_wealth_amount"`
	TaxableObjectSalesValue      string   `json:"taxable_object_sales_value"`
	Deposito                     string   `json:"deposito"`
	Currency                     string   `json:"currency"`
	InvestmentGoals              string   `json:"investment_goals"`
	InvestmentExperience         string   `json:"investment_experience"`
	BankList                     BankList `gorm:"type:jsonb" json:"bank_list"` // Column for the array of obj
	CurrencyRate                 float64  `json:"currency_rate"`
	ProductServicePlatform       string   `json:"product_service_platform"`
	DomAddress                   string   `json:"dom_address"`
	IdentityAddress              string   `json:"identity_address"`
	AccountType                  string   `json:"account_type"` // e.g., Personal, Joint, etc.
}

type UserEmergencyContact struct {
	ID                          string `gorm:"primary_key" json:"userid"`
	EmergencyContactName        string `json:"emergency_contact_name"`
	EmergencyContactCountry     string `json:"emergency_contact_country"`
	EmergencyContactProvince    string `json:"emergency_contact_province"`
	EmergencyContactCity        string `json:"emergency_contact_city"`
	EmergencyContactDistrict    string `json:"emergency_contact_district"`
	EmergencyContactSubDistrict string `json:"emergency_contact_subdistrict"`
	EmergencyContactAddress     string `json:"emergency_contact_address"`
	EmergencyContactPhone       string `json:"emergency_contact_phone"`
	EmergencyContactRelation    string `json:"emergency_contact_relation"`

	base.BaseModel
}

type BackOfficeUserLeads struct {
	Userid      string `json:"userid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobile_phone"`
}

// AfterCreate GORM hook to simulate a trigger
func (u *Users) AfterCreate(tx *gorm.DB) (err error) {
	if err := tx.Save(&UserProfile{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserAddress{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserEmployment{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserFinance{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserEmergencyContact{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	accountID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if err := tx.Save(&Account{
		ID:                common.UUIDNormalizer(accountID),
		Type:              enums.AccountTypeDemo,
		ApprovalStatus:    enums.AccountApprovalStatusApproved,
		UserID:            u.ID,
		IsDemo:            true,
		Balance:           1000,
		Equity:            1000,
		MetaLoginId:       u.MetaLoginId,
		MetaLoginPassword: u.MetaLoginPassword,

		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}
func (u *Users) AfterUpdateUserByMobilePhone(tx *gorm.DB) (err error) {
	if err := tx.Save(&UserProfile{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserAddress{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserEmployment{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserFinance{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	if err := tx.Save(&UserEmergencyContact{
		ID: u.ID,
	}).Error; err != nil {
		return err
	}
	accountID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if err := tx.Save(&Account{
		ID:                common.UUIDNormalizer(accountID),
		Type:              enums.AccountTypeDemo,
		ApprovalStatus:    enums.AccountApprovalStatusApproved,
		UserID:            u.ID,
		IsDemo:            true,
		Balance:           1000,
		Equity:            1000,
		MetaLoginId:       u.MetaLoginId,
		MetaLoginPassword: u.MetaLoginPassword,

		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

type UserLoginSessionRequest struct {
	MobilePhone string `json:"mobilephone" binding:"required"`
	Pin         string `json:"pin" validate:"required,exact=6"`
	Device      string `json:"device"`
	DeviceId    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
	DeviceImei  string `json:"device_imei"`
	DeviceOs    string `json:"device_os"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	SessionId   string `json:"session_id"`
}

type UserSessionRequest struct {
	DeviceId  string `json:"device_id"`
	SessionId string `json:"session_id"`
}

type UpdatePinRequest struct {
	MobilePhone string `json:"mobilephone" binding:"required"`
	Pin         string `json:"pin" validate:"required,exact=6"`
}

type UserLoginResponse struct {
	IsRegistered bool `json:"is_registered"`
}

type UserLoginSessionResponse struct {
	Name             string    `json:"customer_name"`
	Email            string    `json:"email"`
	AccessToken      string    `json:"access_token"`
	ExpirationString time.Time `json:"expiration_string"`
	Expiration       int64     `json:"expiration"`
	// SessionId string `json:"session_id"`
}

type UserApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}
