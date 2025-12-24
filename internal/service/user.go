package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/duckvoid/yago-mart/internal/domain/balance"
	userdomain "github.com/duckvoid/yago-mart/internal/domain/user"
)

type UserService struct {
	repo       userdomain.Repository
	logger     *slog.Logger
	balanceSvc *BalanceService
}

func NewUserService(repo userdomain.Repository, balanceSvc *BalanceService, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, balanceSvc: balanceSvc, logger: logger}
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

func (u *UserService) GetBalance(ctx context.Context, username string) (*balance.Entity, error) {
	return u.balanceSvc.Get(ctx, username)
}

func (u *UserService) Accrual(ctx context.Context, username string, value float64) error {
	return u.balanceSvc.Accrual(ctx, username, value)
}

func (u *UserService) Withdrawal(ctx context.Context, username string, value float64) error {
	return u.balanceSvc.Withdrawal(ctx, username, value)
}
