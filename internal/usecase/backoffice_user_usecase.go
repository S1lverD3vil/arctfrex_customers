package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"arctfrex-customers/internal/auth"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
)

type BackofficeUserUsecase interface {
	Register(backofficeUser *model.BackofficeUsers) error
	LoginSession(backofficeUser *model.BackofficeUsers) (*model.BackofficeUserLoginSessionResponse, error)
	All() (*[]model.BackofficeUsers, error)
	AllUsersByRoleId(roleId string) ([]model.BackofficeUsers, error)
	Subordinate(userId string) (*[]model.BackofficeUsers, error)
}

type backofficeUserUsecase struct {
	backofficeUserRepository repository.BackofficeUserRepository
	tokenService             auth.TokenService
}

func NewBackofficeUsecase(
	bur repository.BackofficeUserRepository,
	ts auth.TokenService,
) *backofficeUserUsecase {
	return &backofficeUserUsecase{
		backofficeUserRepository: bur,
		tokenService:             ts,
	}
}

func (buu *backofficeUserUsecase) Register(backofficeUser *model.BackofficeUsers) error {
	userdb, _ := buu.backofficeUserRepository.GetUserByEmail(backofficeUser.Email)
	if userdb != nil {
		return errors.New("email or phone number already used")
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(backofficeUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	backofficeUser.ID = strings.ReplaceAll(newUUID.String(), "-", "")
	backofficeUser.Password = string(hashedPassword)
	backofficeUser.ReferralCode = common.GenerateShortKSUID()
	backofficeUser.IsActive = true
	fmt.Println(backofficeUser.JobPosition.String())

	return buu.backofficeUserRepository.Create(backofficeUser)
}

func (buu *backofficeUserUsecase) LoginSession(backofficeUser *model.BackofficeUsers) (*model.BackofficeUserLoginSessionResponse, error) {
	userdb, err := buu.backofficeUserRepository.GetUserByEmail(backofficeUser.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userdb.Password), []byte(backofficeUser.Password))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	generatedToken, expirationTime, err := buu.tokenService.GenerateToken(userdb.ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	userdb.SessionId = generatedToken
	if err := buu.backofficeUserRepository.Update(userdb); err != nil {
		return nil, err
	}

	return &model.BackofficeUserLoginSessionResponse{
		ID:           userdb.ID,
		Name:         userdb.Name,
		Email:        userdb.Email,
		AccessToken:  generatedToken,
		Expiration:   int64(time.Until(expirationTime).Seconds()),
		RoleId:       userdb.RoleId,
		ReferralCode: userdb.ReferralCode,
	}, nil
}

func (buh *backofficeUserUsecase) All() (*[]model.BackofficeUsers, error) {
	userdb, err := buh.backofficeUserRepository.GetActiveUsers()
	if err != nil {
		return nil, err
	}

	return userdb, nil
}

func (buh *backofficeUserUsecase) AllUsersByRoleId(roleId string) ([]model.BackofficeUsers, error) {
	userdb, err := buh.backofficeUserRepository.GetActiveUsersByRoleId(roleId)
	if err != nil {
		return nil, err
	}

	return userdb, nil
}

func (buh *backofficeUserUsecase) Subordinate(userId string) (*[]model.BackofficeUsers, error) {
	userdb, err := buh.backofficeUserRepository.GetActiveSubordinate(userId)
	if err != nil {
		return nil, err
	}

	return userdb, nil
}
