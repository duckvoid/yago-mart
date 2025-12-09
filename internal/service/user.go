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

func (u *UserService) All() ([]*userdomain.User, error) {
	return u.repo.All()
}

func (u *UserService) Get(login string, password string) (*userdomain.User, error) {
	return u.repo.Get(login, password)
}

func (u *UserService) Create(username string, password string) error {
	user := &userdomain.User{
		Name:     username,
		Password: password,
	}
	return u.repo.Create(user)
}
