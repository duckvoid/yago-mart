package service

import (
	"strconv"

	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
)

type BalanceService struct {
	repo     balancedomain.Repository
	orderSvc *OrderService
}

func NewBalanceService(repo balancedomain.Repository, orderSvc *OrderService) *BalanceService {
	return &BalanceService{
		repo:     repo,
		orderSvc: orderSvc,
	}
}

func (b *BalanceService) Get(username string) (*balancedomain.Entity, error) {
	balance, err := b.repo.Get(username)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (b *BalanceService) Accrual(username string, value float64) error {
	return b.repo.Accrual(username, value)
}

func (b *BalanceService) Withdrawal(username string, orderID string, value float64) error {
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return err
	}

	if _, err = b.orderSvc.Get(id); err != nil {
		return err
	}

	return b.repo.Withdrawal(username, value)
}
