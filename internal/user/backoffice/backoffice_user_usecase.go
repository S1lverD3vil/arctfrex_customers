package user

import (
	"arctfrex-customers/internal/auth"
	"arctfrex-customers/internal/common"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type BackofficeUserUsecase interface {
	Register(backofficeUser *BackofficeUsers) error
	LoginSession(backofficeUser *BackofficeUsers) (*BackofficeUserLoginSessionResponse, error)
	All() (*[]BackofficeUsers, error)
}

type backofficeUserUsecase struct {
	backofficeUserRepository BackofficeUserRepository
	tokenService             auth.TokenService
}

func NewBackofficeUsecase(
	bur BackofficeUserRepository,
	ts auth.TokenService,
) *backofficeUserUsecase {
	return &backofficeUserUsecase{
		backofficeUserRepository: bur,
		tokenService:             ts,
	}
}

func (buu *backofficeUserUsecase) Register(backofficeUser *BackofficeUsers) error {

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
	backofficeUser.RoleId = backofficeUser.RoleIdType.String()
	backofficeUser.Password = string(hashedPassword)
	backofficeUser.ReferralCode = common.GenerateShortKSUID()
	backofficeUser.IsActive = true
	fmt.Println(backofficeUser.JobPosition.String())

	return buu.backofficeUserRepository.Create(backofficeUser)
}

func (buu *backofficeUserUsecase) LoginSession(backofficeUser *BackofficeUsers) (*BackofficeUserLoginSessionResponse, error) {
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

	return &BackofficeUserLoginSessionResponse{
		ID:           userdb.ID,
		Name:         userdb.Name,
		Email:        userdb.Email,
		AccessToken:  generatedToken,
		Expiration:   int64(time.Until(expirationTime).Seconds()),
		RoleId:       userdb.RoleId,
		ReferralCode: userdb.ReferralCode,
	}, nil
}

func (buh *backofficeUserUsecase) All() (*[]BackofficeUsers, error) {
	userdb, err := buh.backofficeUserRepository.GetActiveUsers()
	if err != nil {
		return nil, err
	}

	return userdb, nil

}
