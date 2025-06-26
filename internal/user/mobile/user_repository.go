package user

import (
	"fmt"

	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
)

// userRepository struct implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (ur *userRepository) Create(user *Users) error {
	return ur.db.Create(user).Error
}

// Create inserts a new user into the database
func (ur *userRepository) Save(user *Users) error {
	return ur.db.Save(user).Error
}

func (ur *userRepository) GetUserByEmail(email string) (*Users, error) {
	var user Users
	if err := ur.db.Where(&Users{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetActiveUserByUserId(userId string) (*Users, error) {
	var user Users
	queryParams := Users{
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

func (ur *userRepository) GetActiveUserByUserIdSessionId(userId, sessionId string) (*Users, error) {
	var user Users
	queryParams := Users{
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

func (ur *userRepository) GetActiveUserProfileByUserId(userId string) (*UserProfile, error) {
	var userProfile UserProfile
	queryParams := UserProfile{
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

func (ur *userRepository) GetActiveUserProfileDetailByUserID(userID string) (*UserProfileDetail, error) {
	var userProfileDetail UserProfileDetail

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
			users.fax_number,
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
		Scan(&userProfileDetail).Error; err != nil {

		return nil, err
	}

	return &userProfileDetail, nil
}

func (ur *userRepository) GetActiveUserAddressByUserId(userId string) (*UserAddress, error) {
	var userAddress UserAddress
	queryParams := UserAddress{
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

func (ur *userRepository) GetActiveUserEmploymentByUserId(userId string) (*UserEmployment, error) {
	var userEmployment UserEmployment
	queryParams := UserEmployment{
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

func (ur *userRepository) GetActiveUserFinanceByUserId(userId string) (*UserFinance, error) {
	var userFinance UserFinance
	queryParams := UserFinance{
		ID: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := ur.db.Where(&queryParams).First(&userFinance).Error; err != nil {
		return nil, err
	}

	return &userFinance, nil
}

func (ur *userRepository) GetActiveUserEmergencyContactByUserId(userId string) (*UserEmergencyContact, error) {
	var userEmergencyContact UserEmergencyContact
	queryParams := UserEmergencyContact{
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

func (ur *userRepository) GetActiveUserByMobilePhone(mobilePhone string) (*Users, error) {
	var user Users
	queryParams := Users{
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

func (ur *userRepository) GetUserByEmailAndMobilePhone(email, mobilePhone string) (*Users, error) {
	var user Users
	queryParams := Users{
		Email:       email,
		MobilePhone: mobilePhone,
	}
	if err := ur.db.Where(&queryParams).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByMobilePhone(mobilePhone string) (*Users, error) {
	var user Users
	if err := ur.db.Where(&Users{MobilePhone: mobilePhone}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByEmailOrMobilePhone(email, mobilePhone string) (*Users, error) {
	var user Users
	if err := ur.db.Where(&Users{Email: email}).Or(&Users{MobilePhone: mobilePhone}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetActiveUserLeads() (*[]BackOfficeUserLeads, error) {
	var backOfficeUserLeads []BackOfficeUserLeads
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

func (ur *userRepository) Update(user *Users) error {
	return ur.db.Updates(user).Error
}

func (ur *userRepository) UpdateUserByMobilePhone(user *Users) error {
	if err := ur.db.Model(&Users{}).Where("mobile_phone = ?", user.MobilePhone).Updates(user).Error; err != nil {
		return err
	}
	if user.Pin != common.STRING_EMPTY {
		user.AfterCreate(ur.db)
	}

	return nil
}

func (ur *userRepository) UpdateUserWatchlist(user *Users) error {
	return ur.db.Select("Watchlist").Updates(user).Error
}

func (ur *userRepository) UpdateProfile(userProfile *UserProfile) error {
	return ur.db.Updates(userProfile).Error
}

func (ur *userRepository) UpdateLogoutSession(user *Users) error {
	return ur.db.Select("SessionId").Updates(user).Error
}

func (ur *userRepository) UpdateDeleteUser(user *Users) error {
	return ur.db.Select("IsActive").Updates(user).Error
}

func (ur *userRepository) UpdateProfileKtpPhoto(userProfile *UserProfile) error {
	return ur.db.Select("KtpPhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileSelfiePhoto(userProfile *UserProfile) error {
	return ur.db.Select("SelfiePhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileNpwpPhoto(userProfile *UserProfile) error {
	return ur.db.Select("NpwpPhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileAdditionalDocumentPhoto(userProfile *UserProfile) error {
	return ur.db.Select("AdditionalDocumentPhoto").Updates(userProfile).Error
}

func (ur *userRepository) UpdateProfileDeclarationVideo(userProfile *UserProfile) error {
	return ur.db.Select("DeclarationVideo").Updates(userProfile).Error
}

func (ur *userRepository) UpdateAddress(userAddress *UserAddress) error {
	return ur.db.Updates(userAddress).Error
}

func (ur *userRepository) UpdateEmployment(userEmployment *UserEmployment) error {
	return ur.db.Updates(userEmployment).Error
}

func (ur *userRepository) UpdateFinance(userFinance *UserFinance) error {
	return ur.db.Updates(userFinance).Error
}

func (ur *userRepository) UpdateEmergencyContact(userEmergencyContact *UserEmergencyContact) error {
	return ur.db.Updates(userEmergencyContact).Error
}
