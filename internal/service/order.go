package service

import (
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

func (o *OrderService) Create(username string, orderID int) error {

	//accrual := o.accrualSvc.Get(orderID)

	order := &model.Order{
		ID:       orderID,
		Username: username,
		Status:   model.OrderRegistered,
		//Accrual: accrual,
	}

	return o.repo.Create(order)
}

func (o *OrderService) UserOrders(username string) ([]*model.Order, error) {
	return o.repo.GetByUser(username)
}

func (o *OrderService) LuhnValidation(number int) bool {

	var digits []int

	for number > 0 {
		digits = append(digits, number%10)
		number /= 10
	}

	sum := 0

	for i := 0; i < len(digits); i++ {
		digit := digits[i]
		if i%2 != 0 {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}
