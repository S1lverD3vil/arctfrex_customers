package repository

import (
	"fmt"

	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
)

type UserRepository interface {
	Create(user *model.Users) error
	Save(user *model.Users) error

	GetUserByEmail(email string) (*model.Users, error)
	GetActiveUserByUserId(userId string) (*model.Users, error)
	GetActiveUserByUserIdSessionId(userId, sessionId string) (*model.Users, error)
	GetActiveUserProfileByUserId(userId string) (*model.UserProfile, error)
	GetActiveUserProfileDetailByUserID(userID string) (*model.UserProfileDetail, error)
	GetActiveUserAddressByUserId(userId string) (*model.UserAddress, error)
	GetActiveUserEmploymentByUserId(userId string) (*model.UserEmployment, error)
	GetActiveUserFinanceByUserId(userId string) (*model.UserFinanceDetail, error)
	GetActiveUserEmergencyContactByUserId(userId string) (*model.UserEmergencyContact, error)
	GetActiveUserByMobilePhone(mobilePhone string) (*model.Users, error)
	GetUserByEmailAndMobilePhone(email, mobilePhone string) (*model.Users, error)
	GetUserByMobilePhone(mobilePhone string) (*model.Users, error)
	GetUserByEmailOrMobilePhone(email, mobilePhone string) (*model.Users, error)
	GetActiveUserLeads() (*[]model.BackOfficeUserLeads, error)

	Update(user *model.Users) error
	UpdateUserByMobilePhone(user *model.Users) error
	UpdateUserWatchlist(user *model.Users) error
	UpdateProfile(userProfile *model.UserProfile) error
	UpdateLogoutSession(user *model.Users) error
	UpdateDeleteUser(user *model.Users) error
	UpdateProfileKtpPhoto(userProfile *model.UserProfile) error
	UpdateProfileSelfiePhoto(userProfile *model.UserProfile) error
	UpdateProfileNpwpPhoto(userProfile *model.UserProfile) error
	UpdateProfileAdditionalDocumentPhoto(userProfile *model.UserProfile) error
	UpdateProfileDeclarationVideo(userProfile *model.UserProfile) error
	UpdateAddress(userAddress *model.UserAddress) error
	UpdateEmployment(userEmployment *model.UserEmployment) error
	UpdateFinance(userFinance *model.UserFinance) error
	UpdateEmergencyContact(userEmergencyContact *model.UserEmergencyContact) error
}

// userRepository struct implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (ur *userRepository) Create(user *model.Users) error {
	return ur.db.Create(user).Error
}

// Create inserts a new user into the database
func (ur *userRepository) Save(user *model.Users) error {
	return ur.db.Save(user).Error
}

func (ur *userRepository) GetUserByEmail(email string) (*model.Users, error) {
	var user model.Users
	if err := ur.db.Where(&model.Users{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetActiveUserByUserId(userId string) (*model.Users, error) {
	var user model.Users
	queryParams := model.Users{
		ID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetActiveUserByUserIdSessionId(userId, sessionId string) (*model.Users, error) {
	var user model.Users
	queryParams := model.Users{
		ID:        userId,
		SessionId: sessionId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetActiveUserProfileByUserId(userId string) (*model.UserProfile, error) {
	var userProfile model.UserProfile
	queryParams := model.UserProfile{
		ID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&userProfile).Error; err != nil {
		return nil, err
	}

	return &userProfile, nil
}

func (ur *userRepository) GetActiveUserProfileDetailByUserID(userID string) (*model.UserProfileDetail, error) {
	var userProfileDetail model.UserProfileDetail

	if err := ur.db.Table("users").
		Joins(`
			LEFT JOIN user_profiles 
				ON users.id = user_profiles.id 
				AND user_profiles.is_active = ?`,
			true,
		).
		Joins(`			
			LEFT JOIN user_addresses 
				ON users.id = user_addresses.id 
				AND user_addresses.is_active = ?`,
			true,
		).
		Select(`
			users.id as user_id, 
			users.name as full_name, 
			users.mobile_phone as mobile_phone,
			users.home_phone,
			users.fax_phone,
			user_addresses.dom_postal_code,
			user_profiles.identity_type,
			user_profiles.gender,
			user_profiles.place_of_birth,
			user_profiles.marital_status,
			user_profiles.date_of_birth,
			user_profiles.ktp_number,
			user_profiles.ktp_photo,
			user_profiles.selfie_photo,
			user_profiles.nationality,
			user_profiles.npwp_number,
			user_profiles.npwp_photo,
			user_profiles.additional_document_photo,
			user_profiles.declaration_video,
			user_profiles.mother_maiden,
			user_profiles.created_at,
			user_profiles.modified_at
		`).
		Where("users.id =?", userID).
		Where("users.is_active = ?", true).
		Scan(&userProfileDetail).Error; err != nil {

		return nil, err
	}

	return &userProfileDetail, nil
}

func (ur *userRepository) GetActiveUserAddressByUserId(userId string) (*model.UserAddress, error) {
	var userAddress model.UserAddress
	queryParams := model.UserAddress{
		ID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&userAddress).Error; err != nil {
		return nil, err
	}

	return &userAddress, nil
}

func (ur *userRepository) GetActiveUserEmploymentByUserId(userId string) (*model.UserEmployment, error) {
	var userEmployment model.UserEmployment
	queryParams := model.UserEmployment{
		ID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&userEmployment).Error; err != nil {
		return nil, err
	}

	return &userEmployment, nil
}

func (ur *userRepository) GetActiveUserFinanceByUserId(userId string) (*model.UserFinanceDetail, error) {
	var userFinanceDetail model.UserFinanceDetail

	if err := ur.db.Table("users").
		Joins(`
			LEFT JOIN user_addresses 
				ON users.id = user_addresses.id 
				AND user_addresses.is_active = ?`,
			true,
		).
		Joins(`			
			LEFT JOIN user_finances 
				ON users.id = user_finances.id
				AND user_finances.is_active = ?`,
			true,
		).
		Select(`
			users.id as user_id, 
			user_addresses.dom_address,
			user_addresses.ktp_address AS identity_address,
			user_finances.source_income,
			user_finances.yearly_income_amount,
			user_finances.yearly_additional_income_amount,
			user_finances.estimation_wealth_amount,
			user_finances.taxable_object_sales_value,
			user_finances.deposito,
			user_finances.currency,
			user_finances.investment_goals,
			user_finances.investment_experience,
			user_finances.bank_list,
			user_finances.currency_rate,
			user_finances.product_service_platform,
			user_finances.account_type
		`).
		Where("users.id =?", userId).
		Where("users.is_active = ?", true).
		Scan(&userFinanceDetail).Error; err != nil {
		return nil, err
	}

	return &userFinanceDetail, nil
}

func (ur *userRepository) GetActiveUserEmergencyContactByUserId(userId string) (*model.UserEmergencyContact, error) {
	var userEmergencyContact model.UserEmergencyContact
	queryParams := model.UserEmergencyContact{
		ID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&userEmergencyContact).Error; err != nil {
		return nil, err
	}

	return &userEmergencyContact, nil
}

func (ur *userRepository) GetActiveUserByMobilePhone(mobilePhone string) (*model.Users, error) {
	var user model.Users
	queryParams := model.Users{
		MobilePhone: mobilePhone,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByEmailAndMobilePhone(email, mobilePhone string) (*model.Users, error) {
	var user model.Users
	queryParams := model.Users{
		Email:       email,
		MobilePhone: mobilePhone,
	}
	if err := ur.db.Where(&queryParams).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByMobilePhone(mobilePhone string) (*model.Users, error) {
	var user model.Users
	if err := ur.db.Where(&model.Users{MobilePhone: mobilePhone}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByEmailOrMobilePhone(email, mobilePhone string) (*model.Users, error) {
	var user model.Users
	if err := ur.db.Where(&model.Users{Email: email}).Or(&model.Users{MobilePhone: mobilePhone}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetActiveUserLeads() (*[]model.BackOfficeUserLeads, error) {
	var backOfficeUserLeads []model.BackOfficeUserLeads
	if err := ur.db.Table("users").
		Joins(`
			LEFT JOIN accounts 
				ON users.id = accounts.user_id
				AND accounts.type = ?
				AND users.is_active = ?
			AND accounts.is_active = ?`,
			enums.AccountTypeReal,
			true,
			true,
		).
		Select(`
			users.id as userid, 
			users.name as name, 
			users.email,
			users.mobile_phone as mobile_phone
		`).
		Where("accounts.user_id is null").
		Scan(&backOfficeUserLeads).Error; err != nil {

		return nil, err
	}
	fmt.Println(&backOfficeUserLeads)
	return &backOfficeUserLeads, nil
}

func (ur *userRepository) Update(user *model.Users) error {
	return ur.db.Updates(user).Error
}

func (ur *userRepository) UpdateUserByMobilePhone(user *model.Users) error {
	if err := ur.db.Model(&model.Users{}).Where("mobile_phone = ?", user.MobilePhone).Updates(user).Error; err != nil {
		return err
	}
	if user.Pin != common.STRING_EMPTY {
		user.AfterCreate(ur.db)
	}

	return nil
}

func (ur *userRepository) UpdateUserWatchlist(user *model.Users) error {
	return ur.db.Select("Watchlist").Updates(user).Error
}

func (ur *userRepository) UpdateProfile(userProfile *model.UserProfile) error {
	return ur.db.Updates(userProfile).Error
}

func (ur *userRepository) UpdateLogoutSession(user *model.Users) error {
	return ur.db.Select("SessionId").Updates(user).Error
}

func (ur *userRepository) UpdateDeleteUser(user *model.Users) error {
	return ur.db.Select("IsActive").Updates(user).Error
}

func (ur *userRepository) UpdateProfileKtpPhoto(userProfile *model.UserProfile) error {
	return ur.db.Select("KtpPhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileSelfiePhoto(userProfile *model.UserProfile) error {
	return ur.db.Select("SelfiePhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileNpwpPhoto(userProfile *model.UserProfile) error {
	return ur.db.Select("NpwpPhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileAdditionalDocumentPhoto(userProfile *model.UserProfile) error {
	return ur.db.Select("AdditionalDocumentPhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileDeclarationVideo(userProfile *model.UserProfile) error {
	return ur.db.Select("DeclarationVideo").Updates(userProfile).Error
}

func (ur *userRepository) UpdateAddress(userAddress *model.UserAddress) error {
	return ur.db.Updates(userAddress).Error
}

func (ur *userRepository) UpdateEmployment(userEmployment *model.UserEmployment) error {
	return ur.db.Updates(userEmployment).Error
}

func (ur *userRepository) UpdateFinance(userFinance *model.UserFinance) error {
	return ur.db.Updates(userFinance).Error
}

func (ur *userRepository) UpdateEmergencyContact(userEmergencyContact *model.UserEmergencyContact) error {
	return ur.db.Updates(userEmergencyContact).Error
}
