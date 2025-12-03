package service

import "github.com/duckvoid/yago-mart/internal/model"

type UserRepository interface {
	All() []*model.User
	Get(id int64) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int64) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(login string, password string) error {
	user := &model.User{
		Login:    login,
		Password: password,
	}
	return s.repo.Create(user)
}
