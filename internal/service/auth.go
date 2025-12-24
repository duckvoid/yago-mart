package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var singingKey []byte

type AuthService struct {
	userSvc *UserService
	logger  *slog.Logger
}

func NewAuthService(signingKey string, userSvc *UserService, logger *slog.Logger) *AuthService {
	singingKey = []byte(signingKey)
	return &AuthService{
		userSvc: userSvc,
		logger:  logger,
	}
}

func (a *AuthService) Register(ctx context.Context, username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("failed to generate password")
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
		a.logger.Error("User passwords do not match", "username", username)
		return "", err
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(singingKey)
	if err != nil {
		a.logger.Error("Failed to sign token", "username", username)
		return "", err
	}

	return tokenString, nil
}

func AuthToken(authToken string) (string, error) {
	token, err := jwt.ParseWithClaims(authToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return singingKey, nil
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
