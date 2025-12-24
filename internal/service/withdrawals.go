package service

import (
	"context"
	"errors"
	"log/slog"

	withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"
)

type WithdrawalsService struct {
	repo   withdrawalsdomain.Repository
	logger *slog.Logger
}

func NewWithdrawalsService(repo withdrawalsdomain.Repository, logger *slog.Logger) *WithdrawalsService {
	return &WithdrawalsService{repo: repo, logger: logger}
}

func (w *WithdrawalsService) UserWithdrawals(ctx context.Context, username string) ([]*withdrawalsdomain.Entity, error) {
	withdrawals, err := w.repo.GetByUser(ctx, username)
	if err != nil {
		if errors.Is(err, withdrawalsdomain.ErrNotFound) {
			w.logger.Error("user withdrawals not found", "username", username)
			return nil, withdrawalsdomain.ErrNotFound
		}

		return nil, err
	}
	return withdrawals, nil
}

func (w *WithdrawalsService) Create(ctx context.Context, username string, orderID int, sum float64) error {
	return w.repo.Create(ctx, &withdrawalsdomain.Entity{
		Username: username,
		OrderID:  orderID,
		Sum:      sum,
	})
}
