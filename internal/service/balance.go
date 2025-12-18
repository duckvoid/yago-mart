package service

import (
	"context"
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

func (b *BalanceService) Get(ctx context.Context, username string) (*balancedomain.Entity, error) {
	balance, err := b.repo.Get(ctx, username)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (b *BalanceService) Accrual(ctx context.Context, username string, value float64) error {
	return b.repo.Accrual(ctx, username, value)
}

func (b *BalanceService) Withdrawal(ctx context.Context, username string, orderID string, value float64) error {
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return err
	}

	if _, err = b.orderSvc.Get(ctx, id); err != nil {
		return err
	}

	return b.repo.Withdrawal(ctx, username, value)
}
