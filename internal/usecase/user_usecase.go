package usecase

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"arctfrex-customers/internal/api"
	"arctfrex-customers/internal/auth"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
)

type UserUsecase interface {
	Register(user *model.Users) error
	Check(mobilePhone string) (*model.Users, error)
	LoginSession(user *model.UserLoginSessionRequest) (*model.UserLoginSessionResponse, error)
	// LoginSession(mobilePhone, password, deviceId string) (*UserLoginSessionResponse, error)
	Session(user *model.Users) (*model.Users, error)
	LogoutSession(user *model.Users) error
	Delete(user *model.Users) error
	UpdatePin(mobilePhone, pin string) error
	UpdateProfile(userID string, userProfile *model.UserProfile) error
	GetProfile(userID string) (*model.UserProfileDetail, error)
	UpdateAddress(userID string, userAddress *model.UserAddress) error
	GetAddress(userID string) (*model.UserAddress, error)
	UpdateEmployment(userID string, userEmployment *model.UserEmployment) error
	GetEmployment(userID string) (*model.UserEmploymentDetail, error)
	UpdateFinance(userID string, userFinance *model.UserFinance) error
	GetFinance(userID string) (*model.UserFinanceDetail, error)
	UpdateEmergencyContact(userID string, userEmegencyContact *model.UserEmergencyContact) error
	GetEmergencyContact(userID string) (*model.UserEmergencyContact, error)
	BackOfficeGetLeads() (*[]model.BackOfficeUserLeads, error)
}

type userUsecase struct {
	userRepository repository.UserRepository
	tokenService   auth.TokenService
	userApiclient  api.UserApiclient
}

func NewUserUseCase(
	ur repository.UserRepository,
	ts auth.TokenService,
	ua api.UserApiclient,
) *userUsecase {
	return &userUsecase{
		userRepository: ur,
		tokenService:   ts,
		userApiclient:  ua,
	}
}

func (uu *userUsecase) Register(user *model.Users) error {
	user.Email = strings.ToLower(user.Email)
	userdb, _ := uu.userRepository.GetUserByMobilePhone(user.MobilePhone)
	// userdb, _ := uu.userRepository.GetUserByEmailAndMobilePhone(user.Email, user.MobilePhone)

	//Already active user
	if userdb != nil && userdb.IsActive {
		return errors.New("email or phone number already used")
	}

	// deletedUser := userdb != nil && userdb.Pin != common.STRING_EMPTY
	// fmt.Printf("%+v", userdb)
	// fmt.Println(userdb.Pin)
	// fmt.Println(userdb.Pin != common.STRING_EMPTY)
	// fmt.Println(deletedUser)

	// log.Println(userdb.Pin == common.STRING_EMPTY)

	// if userdb != nil && userdb.MetaLoginId != 0 {
	// 	user.MetaLoginId = userdb.MetaLoginId
	// 	user.MetaLoginPassword = userdb.MetaLoginPassword
	// }
	//Not active or Deleted user
	// if !deletedUser {
	if userdb != nil && userdb.Pin == common.STRING_EMPTY {
		// if userdb != nil && userdb.Pin == common.STRING_EMPTY {
		user.ID = userdb.ID
		user.MetaLoginId = userdb.MetaLoginId
		user.MetaLoginPassword = userdb.MetaLoginPassword
		// user.MetaLoginId = 0
		// user.MetaLoginPassword = common.STRING_EMPTY

		// if userdb.MetaLoginId != 0 {
		// 	user.MetaLoginId = userdb.MetaLoginId
		// 	user.MetaLoginPassword = userdb.MetaLoginPassword
		// }

		// if userdb.MobilePhone != user.MobilePhone || userdb.Email != user.Email {
		// 	userdb = nil
		// }

		// if userdb.MobilePhone == user.MobilePhone && userdb.Email == user.Email {
		// 	user.MetaLoginId = 0
		// 	user.MetaLoginPassword = common.STRING_EMPTY
		// 	// fmt.Println("masuk sini")
		// 	// fmt.Println(user)
		// 	// fmt.Println(user.MetaLoginId)
		// 	// fmt.Println(user.MetaLoginId == 0)
		// }
		// fmt.Println("masuk sini")
		// fmt.Println(userdb)
	}

	var clientAdd model.ClientAdd
	//Registering demo account to mt5
	if user.MetaLoginId == 0 {
		securedPassword, err := common.GenerateSecurePassword()
		if err != nil {
			return err
		}
		clientAdd = model.ClientAdd{
			Name:     user.Name,
			Password: securedPassword,
			Group:    "demo\\PKB\\B-USD-SFL-MAR-C5-SWAP",
			Leverage: 100,
			Email:    user.Email,
		}
		// fmt.Printf("%+v", clientAdd)

		clientAddData, err := uu.userApiclient.ClientAdd(clientAdd)
		if err != nil {
			return err
		}

		demoAccountTopUp := model.DemoAccountTopUp{
			Login:  clientAddData.Login,
			Amount: 1000,
		}
		demoAccountTopUpData, err := uu.userApiclient.DemoAccountTopUp(demoAccountTopUp)
		if err != nil {
			log.Println(demoAccountTopUpData.Result)
		}

		user.MetaLoginPassword = clientAdd.Password
		user.MetaLoginId = clientAddData.Login
	}

	// if deletedUser {
	if userdb != nil && userdb.Pin != common.STRING_EMPTY {
		user.Pin = common.GenerateShortKSUID()

		newUUID, err := uuid.NewUUID()
		if err != nil {
			log.Println(err)
			return err
		}
		user.ID = common.UUIDNormalizer(newUUID)
	}

	if userdb != nil {
		// user.ID = userdb.ID

		return uu.userRepository.UpdateUserByMobilePhone(user)
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return err
	}
	user.ID = common.UUIDNormalizer(newUUID)

	return uu.userRepository.Create(user)
}

// 200 -> input pin
// 400 -> create pin
// 402 -> signup
func (uu *userUsecase) Check(mobilePhone string) (*model.Users, error) {
	user, err := uu.userRepository.GetActiveUserByMobilePhone(mobilePhone)
	// user, err := uu.userRepository.GetUserByMobilePhone(mobilePhone)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *userUsecase) LoginSession(user *model.UserLoginSessionRequest) (*model.UserLoginSessionResponse, error) {
	// LoginSession(mobilePhone, pin, deviceId string) (*UserLoginSessionResponse, error) {
	userdb, err := uu.userRepository.GetActiveUserByMobilePhone(user.MobilePhone)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userdb.Pin), []byte(user.Pin))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	generatedToken, expirationTime, err := uu.tokenService.GenerateToken(userdb.ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	userdb.SessionId = generatedToken
	userdb.Device = user.Device
	userdb.DeviceOs = user.DeviceOs
	userdb.DeviceId = user.DeviceId
	userdb.DeviceName = user.DeviceName
	userdb.DeviceImei = user.DeviceImei
	userdb.Latitude = user.Latitude
	userdb.Longitude = user.Longitude
	userdb.SessionExpiration = time.Now().Add(48 * time.Hour)
	if err := uu.userRepository.Update(userdb); err != nil {
		return nil, err
	}

	return &model.UserLoginSessionResponse{
		Name:             userdb.Name,
		Email:            userdb.Email,
		AccessToken:      generatedToken,
		ExpirationString: expirationTime,
		Expiration:       int64(time.Until(expirationTime).Seconds()),
	}, err
}

func (uu *userUsecase) Session(user *model.Users) (*model.Users, error) {
	//	fmt.Printf("User Before: %+v\n", user)
	userdb, err := uu.userRepository.GetActiveUserByUserIdSessionId(user.ID, user.SessionId)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	if user.DeviceId != userdb.DeviceId {
		userdb.SessionId = common.STRING_EMPTY
		uu.userRepository.UpdateLogoutSession(userdb)
		return nil, errors.New("unauthorized")
	}

	//fmt.Printf("User After: %+v\n", user)
	userdb.SessionExpiration = time.Now().Add(48 * time.Hour)
	uu.userRepository.Update(userdb)

	return userdb, nil
}

func (uu *userUsecase) LogoutSession(user *model.Users) error {
	user, err := uu.userRepository.GetActiveUserByUserIdSessionId(user.ID, user.SessionId)
	if err != nil {
		return errors.New("unauthorized")
	}

	user.SessionId = common.STRING_EMPTY

	return uu.userRepository.UpdateLogoutSession(user)
}

func (uu *userUsecase) Delete(user *model.Users) error {
	user, err := uu.userRepository.GetActiveUserByUserId(user.ID)
	if err != nil {
		return errors.New("unauthorized")
	}

	user.IsActive = false

	return uu.userRepository.UpdateDeleteUser(user)
}

func (uu *userUsecase) UpdatePin(mobilePhone, pin string) error {
	user, err := uu.userRepository.GetUserByMobilePhone(mobilePhone)
	if user == nil || err != nil {
		return errors.New("user not found")
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Pin = string(hashedPin)
	user.IsActive = true
	user.BaseModel.ModifiedBy = user.ID

	return uu.userRepository.Update(user)
}

func (uu *userUsecase) UpdateProfile(userID string, userProfile *model.UserProfile) error {
	user, err := uu.userRepository.GetActiveUserByUserId(userID)
	if user == nil || err != nil {
		return errors.New("user not found")
	}
	userProfile.ID = userID
	userProfile.IsActive = true
	userProfile.ModifiedBy = userID

	return uu.userRepository.UpdateProfile(userProfile)
}

func (uu *userUsecase) GetProfile(userID string) (*model.UserProfileDetail, error) {

	userProfile, err := uu.userRepository.GetActiveUserProfileDetailByUserID(userID)
	if userProfile == nil || err != nil {
		return nil, errors.New("user not found")
	}

	return userProfile, nil
}

func (uu *userUsecase) UpdateAddress(userID string, userAddress *model.UserAddress) error {
	user, err := uu.userRepository.GetActiveUserByUserId(userID)
	if user == nil || err != nil {
		return errors.New("user not found")
	}
	userAddress.ID = userID
	userAddress.IsActive = true
	userAddress.ModifiedBy = user.ID

	return uu.userRepository.UpdateAddress(userAddress)
}

func (uu *userUsecase) GetAddress(userID string) (*model.UserAddress, error) {
	userAddress, err := uu.userRepository.GetActiveUserAddressByUserId(userID)
	if userAddress == nil || err != nil {
		return nil, errors.New("user not found")
	}

	return userAddress, nil
}

func (uu *userUsecase) UpdateEmployment(userID string, userEmployment *model.UserEmployment) error {
	user, err := uu.userRepository.GetActiveUserByUserId(userID)
	if user == nil || err != nil {
		return errors.New("user not found")
	}
	userEmployment.ID = userID
	userEmployment.IsActive = true
	userEmployment.ModifiedBy = user.ID

	return uu.userRepository.UpdateEmployment(userEmployment)
}

func (uu *userUsecase) GetEmployment(userID string) (*model.UserEmploymentDetail, error) {
	userEmployment, err := uu.userRepository.GetActiveUserEmploymentByUserId(userID)
	if userEmployment == nil || err != nil {
		return nil, errors.New("user not found")
	}

	data := &model.UserEmploymentDetail{
		UserID:            userEmployment.ID,
		CompanyName:       userEmployment.CompanyName,
		CompanyAddress:    userEmployment.CompanyAddress,
		CompanyCity:       userEmployment.CompanyCity,
		CompanyPhone:      userEmployment.CompanyPhone,
		Profession:        userEmployment.Profession,
		CompanyPostalCode: userEmployment.CompanyPostalCode,
		WorkingSince:      userEmployment.WorkingSince,
		WorkingField:      userEmployment.WorkingField,
		PreviewJobTitle:   userEmployment.PreviewJobTitle,
		JobTitle:          userEmployment.JobTitle,
	}

	return data, nil
}

func (uu *userUsecase) UpdateFinance(userID string, userFinance *model.UserFinance) error {
	user, err := uu.userRepository.GetActiveUserByUserId(userID)
	if user == nil || err != nil {
		return errors.New("user not found")
	}
	userFinance.ID = userID
	userFinance.IsActive = true
	userFinance.ModifiedBy = user.ID

	return uu.userRepository.UpdateFinance(userFinance)
}

func (uu *userUsecase) GetFinance(userID string) (*model.UserFinanceDetail, error) {
	userFinance, err := uu.userRepository.GetActiveUserFinanceByUserId(userID)
	if userFinance == nil || err != nil {
		return nil, errors.New("user not found")
	}

	return userFinance, nil
}

func (uu *userUsecase) UpdateEmergencyContact(userID string, userEmegencyContact *model.UserEmergencyContact) error {
	user, err := uu.userRepository.GetActiveUserByUserId(userID)
	if user == nil || err != nil {
		return errors.New("user not found")
	}
	userEmegencyContact.ID = userID
	userEmegencyContact.IsActive = true
	userEmegencyContact.ModifiedBy = user.ID

	return uu.userRepository.UpdateEmergencyContact(userEmegencyContact)
}

func (uu *userUsecase) GetEmergencyContact(userID string) (*model.UserEmergencyContact, error) {
	userEmergencyContact, err := uu.userRepository.GetActiveUserEmergencyContactByUserId(userID)
	if userEmergencyContact == nil || err != nil {
		return nil, errors.New("user not found")
	}

	return userEmergencyContact, nil
}

func (uu *userUsecase) BackOfficeGetLeads() (*[]model.BackOfficeUserLeads, error) {
	backOfficeUserLeads, err := uu.userRepository.GetActiveUserLeads()

	return backOfficeUserLeads, err
}
