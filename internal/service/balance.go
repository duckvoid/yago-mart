package service

import (
	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
)

type BalanceService struct {
	repo balancedomain.Repository
}

func NewBalanceService(repo balancedomain.Repository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

func (b *BalanceService) Get(username string) (float64, float64, error) {
	balance, err := b.repo.Get(username)
	if err != nil {
		return 0, 0, err
	}

	return balance.Current, balance.Withdrawn, nil
}

func (b *BalanceService) Accrual(username string, value float64) error {
	return b.repo.Accrual(username, value)
}

func (b *BalanceService) Withdrawal(username string, value float64) error {
	return b.repo.Withdrawal(username, value)
}
