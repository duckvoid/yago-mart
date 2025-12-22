package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	orderdomain "github.com/duckvoid/yago-mart/internal/domain/order"
)

type OrderService struct {
	repo    orderdomain.Repository
	logger  *slog.Logger
	accrual orderdomain.AccrualClient
}

func NewOrderService(repo orderdomain.Repository, accrual orderdomain.AccrualClient, logger *slog.Logger) *OrderService {
	return &OrderService{repo: repo, accrual: accrual, logger: logger}
}

func (o *OrderService) Create(ctx context.Context, username string, orderID int) error {

	order := &orderdomain.Entity{
		ID:       orderID,
		Username: username,
		Status:   orderdomain.StatusOrderNew,
	}

	err := o.repo.Create(ctx, order)
	if err != nil {
		if errors.Is(err, orderdomain.ErrAlreadyExist) {
			o.logger.Warn("Order already exists", "id", orderID)

			existedOrder, err := o.Get(ctx, orderID)
			if err != nil {
				return err
			}

			if username != existedOrder.Username {
				o.logger.Error("Order already was created by another user", "error", err)
				return orderdomain.ErrCreatedByAnotherUser
			}

		}

		return err
	}

	return nil
}

func (o *OrderService) Get(ctx context.Context, id int) (*orderdomain.Entity, error) {
	order, err := o.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, orderdomain.ErrNotFound) {
			o.logger.Error("Order not found", "id", id)
			return nil, orderdomain.ErrNotFound
		}
		return nil, err
	}
	return order, nil
}

func (o *OrderService) UserOrders(ctx context.Context, username string) ([]*orderdomain.Entity, error) {
	order, err := o.repo.GetByUser(ctx, username)
	if err != nil {
		if errors.Is(err, orderdomain.ErrNotFound) {
			o.logger.Error("Order not found", "username", username)
			return nil, orderdomain.ErrNotFound
		}
		return nil, err
	}
	return order, nil
}

func (o *OrderService) LuhnValidation(orderID int) bool {

	var digits []int

	number := orderID

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

	if sum%10 == 0 {
		return true
	}

	o.logger.Error("Order ID Luhn validation error", slog.Int("number", orderID))

	return false
}

func (o *OrderService) accrualProccess(ctx context.Context, orderID int) error {

	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			accrual, err := o.accrual.GetOrder(timeoutCtx, strconv.Itoa(orderID))
			if err != nil {
				o.logger.Warn("Accrual error", "error", err)
				continue
			}

			switch accrual.Status {
			case orderdomain.StatusAccrualInvalid, orderdomain.StatusAccrualProcessed:

			}

		}
	}

}
