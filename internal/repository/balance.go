package repository

import (
	"context"
	"embed"

	"github.com/duckvoid/yago-mart/internal/model"
	"github.com/jmoiron/sqlx"
)

const BalanceTable = "balance"

//go:embed init_migrations/*init_balance_table.sql
var embedInitBalanceMigration embed.FS

type BalanceRepository struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewBalanceRepository(ctx context.Context, db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{ctx: ctx, db: db}
}

func (b *BalanceRepository) Get(username string) (*model.Balance, error) {
	return nil, nil
}
func (b *BalanceRepository) Accrual(username string, value float64) error {
	return nil
}
func (b *BalanceRepository) Withdrawal(username string, value float64) error {
	return nil
}
