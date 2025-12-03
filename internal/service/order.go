package service

import "github.com/duckvoid/yago-mart/internal/model"

type OrderRepository interface {
	All() []*model.Order
	Get(id int) (*model.Order, error)
	Create(order *model.Order) error
	Update(order *model.Order) error
	Delete(id int) error
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}
