package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenService defines the methods for JWT operations
type TokenService interface {
	GenerateToken(userID string) (string, time.Time, error)
	ValidateToken(token string) (string, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secretKey, issuer string) *jwtService {
	return &jwtService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (s *jwtService) GenerateToken(userID string) (string, time.Time, error) {
	// expirationTime := time.Now().Add(time.Minute * 2)
	expirationTime := time.Now().Add(time.Hour * 72)
	claims := jwt.MapClaims{
		"user_id": userID,
		// "exp":     expirationTime.Format("2006-01-02 15:04:05"),
		"exp": expirationTime.Unix(),
		"iss": s.issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))

	return signedToken, expirationTime, err
}

func (s *jwtService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), nil
	}
	return "", err
}
