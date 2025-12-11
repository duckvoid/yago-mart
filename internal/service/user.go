package service

import (
	userdomain "github.com/duckvoid/yago-mart/internal/domain/user"
)

type UserService struct {
	repo userdomain.Repository
}

func NewUserService(repo userdomain.Repository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) All() ([]*userdomain.Entity, error) {
	return u.repo.All()
}

func (u *UserService) Get(login string) (*userdomain.Entity, error) {
	return u.repo.Get(login)
}

func (u *UserService) Create(username string, password string) error {
	user := &userdomain.Entity{
		Name:     username,
		Password: password,
	}
	return u.repo.Create(user)
}
