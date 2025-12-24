package service

import (
	"context"
	"log/slog"

	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
)

type BalanceService struct {
	repo   balancedomain.Repository
	logger *slog.Logger
}

func NewBalanceService(repo balancedomain.Repository, logger *slog.Logger) *BalanceService {
	return &BalanceService{
		repo:   repo,
		logger: logger,
	}
}

func (b *BalanceService) Get(ctx context.Context, username string) (*balancedomain.Entity, error) {
	return b.repo.Get(ctx, username)
}

func (b *BalanceService) Accrual(ctx context.Context, username string, sum float64) error {
	return b.repo.Accrual(ctx, username, sum)
}

func (b *BalanceService) Withdrawal(ctx context.Context, username string, sum float64) error {
	return b.repo.Withdrawal(ctx, username, sum)
}
