package repository

import (
	"context"
	"database/sql"
	"embed"

	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
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

func (b *BalanceRepository) Get(username string) (*balancedomain.Entity, error) {
	var balance balancedomain.Entity

	row := b.db.QueryRowxContext(b.ctx, `SELECT * FROM balance WHERE user_name = $1`, username)

	if err := row.StructScan(&balance); err != nil {
		return nil, err
	}

	return &balance, nil

}
func (b *BalanceRepository) Accrual(username string, value float64) error {
	return nil
}
func (b *BalanceRepository) Withdrawal(username string, value float64) error {
	tx, err := b.db.BeginTxx(b.ctx, nil)
	if err != nil {
		return err
	}

	var execErr error
	defer func() {
		if execErr != nil {
			_ = tx.Rollback()
		} else {
			execErr = tx.Commit()
		}
	}()

	var res sql.Result
	if res, execErr = tx.ExecContext(b.ctx,
		`UPDATE balance SET current = current - $1 WHERE user_name = $2 AND current >= $1`,
		value, username); err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return balancedomain.ErrInsufficientFunds
	}

	return nil
}
