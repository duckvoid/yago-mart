package service

import (
	"context"
	"errors"
	"log/slog"

	userdomain "github.com/duckvoid/yago-mart/internal/domain/user"
)

type UserService struct {
	repo   userdomain.Repository
	logger *slog.Logger
}

func NewUserService(repo userdomain.Repository, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (u *UserService) All(ctx context.Context) ([]*userdomain.Entity, error) {
	all, err := u.repo.All(ctx)

	if len(all) == 0 {
		u.logger.Error("Failed while fetching all users", "error", err)
		return nil, err
	}

	return all, nil
}

func (u *UserService) Get(ctx context.Context, login string) (*userdomain.Entity, error) {
	user, err := u.repo.Get(ctx, login)
	if err != nil {
		if errors.Is(err, userdomain.ErrNotFound) {
			u.logger.Warn("Failed while fetching user", "login", login, "error", err)
			return nil, userdomain.ErrNotFound
		}

		return nil, err
	}
	return user, nil
}

func (u *UserService) Create(ctx context.Context, username string, password string) error {
	user := &userdomain.Entity{
		Name:     username,
		Password: password,
	}
	return u.repo.Create(ctx, user)
}
