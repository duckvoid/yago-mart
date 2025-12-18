package service

import (
	"context"

	withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"
)

type WithdrawalsService struct {
	repo withdrawalsdomain.Repository
}

func NewWithdrawalsService(repo withdrawalsdomain.Repository) *WithdrawalsService {
	return &WithdrawalsService{repo: repo}
}

func (w *WithdrawalsService) UserWithdrawals(ctx context.Context, username string) ([]*withdrawalsdomain.Entity, error) {
	return w.repo.GetByUser(ctx, username)
}
