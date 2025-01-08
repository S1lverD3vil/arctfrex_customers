package auth

import (
	"os"
	"time"

	"github.com/pkg/errors"
)

type AuthUsecase interface {
	Token(username, password string) (string, time.Time, error)
}

type authUsecase struct {
	tokenService TokenService
}

func NewAuthUsecase(
	ts TokenService,
) *authUsecase {
	return &authUsecase{
		tokenService: ts,
	}
}

func (uc *authUsecase) Token(username, password string) (string, time.Time, error) {

	jwtUsername := os.Getenv("JWT_USERNAME")
	jwtPassword := os.Getenv("JWT_PASSWORD")

	if username != jwtUsername || password != jwtPassword {
		return "", time.Time{}, errors.New("Unauthorized")
	}

	return uc.tokenService.GenerateToken(username)
}
