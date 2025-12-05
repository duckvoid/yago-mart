package service

import (
	"time"

	"github.com/duckvoid/yago-mart/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = "secretKey"

type AuthService struct {
	userSvc *UserService
}

func NewAuthService(userSvc *UserService) *AuthService {
	return &AuthService{userSvc: userSvc}
}

func (a *AuthService) Register(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return a.userSvc.Create(username, string(passwordHash))
}

func (a *AuthService) Login(username, password string) (string, error) {
	user, err := a.userSvc.Get(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Log.Error("CompareHashAndPassword", err.Error())
		return "", err
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.Log.Error("CompareHashAndPassword", err.Error())

		return "", err
	}

	return tokenString, nil

}
