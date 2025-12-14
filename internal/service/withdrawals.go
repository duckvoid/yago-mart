package service

import withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"

type WithdrawalsService struct {
	repo withdrawalsdomain.Repository
}

func NewWithdrawalsService(repo withdrawalsdomain.Repository) *WithdrawalsService {
	return &WithdrawalsService{repo: repo}
}

func (w *WithdrawalsService) UserWithdrawals(username string) ([]*withdrawalsdomain.Entity, error) {
	return w.repo.GetByUser(username)
}
