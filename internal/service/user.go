package service

import "github.com/duckvoid/yago-mart/internal/model"

type UserRepository interface {
	All() ([]*model.User, error)
	Get(username string) (*model.User, error)
	Create(user *model.User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) All() ([]*model.User, error) {
	return u.repo.All()
}

func (u *UserService) Get(login string) (*model.User, error) {
	return u.repo.Get(login)
}

func (u *UserService) Create(username string, password string) error {
	user := &model.User{
		Name:     username,
		Password: password,
	}
	return u.repo.Create(user)
}
