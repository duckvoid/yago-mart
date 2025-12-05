package service

import (
	"time"

	"github.com/duckvoid/yago-mart/internal/model"
)

type OrderRepository interface {
	All() []*model.Order
	Get(id int64) (*model.Order, error)
	GetByUser(username string) ([]*model.Order, error)
	Create(order *model.Order) error
	Update(order *model.Order) error
	Delete(id int64) error
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (o *OrderService) Create(username string, orderID int64) error {

	//accrual := o.accrualSvc.Get(orderID)

	order := &model.Order{
		ID:         orderID,
		Username:   username,
		Status:     "new",
		UploadDate: time.Now(),
		//Accrual: accrual,
	}

	return o.repo.Create(order)
}

func (o *OrderService) UserOrders(username string) ([]*model.Order, error) {
	return o.repo.GetByUser(username)
}
