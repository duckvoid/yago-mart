package service

import (
	"context"

	userdomain "github.com/duckvoid/yago-mart/internal/domain/user"
)

type UserService struct {
	repo userdomain.Repository
}

func NewUserService(repo userdomain.Repository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) All(ctx context.Context) ([]*userdomain.Entity, error) {
	return u.repo.All(ctx)
}

func (u *UserService) Get(ctx context.Context, login string) (*userdomain.Entity, error) {
	return u.repo.Get(ctx, login)
}

func (u *UserService) Create(ctx context.Context, username string, password string) error {
	user := &userdomain.Entity{
		Name:     username,
		Password: password,
	}
	return u.repo.Create(ctx, user)
}
