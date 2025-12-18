package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var singingKey []byte

type AuthService struct {
	userSvc *UserService
}

func NewAuthService(signingKey string, userSvc *UserService) *AuthService {
	singingKey = []byte(signingKey)
	return &AuthService{
		userSvc: userSvc,
	}
}

func (a *AuthService) Register(ctx context.Context, username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return a.userSvc.Create(ctx, username, string(passwordHash))
}

func (a *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.userSvc.Get(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(singingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthToken(authToken string) (string, error) {
	token, err := jwt.ParseWithClaims(authToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(singingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	username, _ := (*claims)["username"].(string)

	expTime, err := claims.GetExpirationTime()
	if err != nil {
		return "", err
	}

	if expTime.Before(time.Now()) {
		return "", errors.New("token is expired")
	}

	return username, nil
}
