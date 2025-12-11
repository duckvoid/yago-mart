package service

import (
	"errors"

	orderdomain "github.com/duckvoid/yago-mart/internal/domain/order"
)

type OrderService struct {
	repo orderdomain.Repository
}

func NewOrderService(repo orderdomain.Repository) *OrderService {
	return &OrderService{repo: repo}
}

func (o *OrderService) Create(username string, orderID int) error {

	//accrual := o.accrualSvc.Get(orderID)

	order := &orderdomain.Entity{
		ID:       orderID,
		Username: username,
		Status:   orderdomain.Registered,
		//Accrual: accrual,
	}

	err := o.repo.Create(order)
	if err != nil {
		if errors.Is(err, orderdomain.ErrAlreadyExist) {
			existedOrder, err := o.Get(orderID)
			if err != nil {
				return err
			}

			if username != existedOrder.Username {
				return orderdomain.ErrCreatedByAnotherUser
			}
		}

		return err
	}

	return nil
}

func (o *OrderService) Get(id int) (*orderdomain.Entity, error) {
	return o.repo.Get(id)
}

func (o *OrderService) UserOrders(username string) ([]*orderdomain.Entity, error) {
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
