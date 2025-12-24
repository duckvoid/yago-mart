package repository

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
	"github.com/jmoiron/sqlx"
)

const BalanceTable = "balance"

//go:embed init_migrations/*init_balance_table.sql
var embedInitBalanceMigration embed.FS

type BalanceRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewBalanceRepository(db *sqlx.DB, logger *slog.Logger) *BalanceRepository {
	return &BalanceRepository{db: db, logger: logger}
}

func (b *BalanceRepository) Get(ctx context.Context, username string) (*balancedomain.Entity, error) {
	var balance balancedomain.Entity

	row := b.db.QueryRowxContext(ctx, `SELECT * FROM balance WHERE user_name = $1`, username)

	if err := row.StructScan(&balance); err != nil {
		b.logger.Error("Failed while getting balance", "user", username, "err", err)
		return nil, err
	}

	return &balance, nil
}

func (b *BalanceRepository) Accrual(ctx context.Context, username string, accrual float64) error {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		b.logger.Error("Failed begin transaction while make withdrawal", "user", username, "err", err)
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

	if _, execErr = tx.ExecContext(ctx,
		`INSERT INTO balance (user_name, current) VALUES ($1, $2) 
				ON CONFLICT(user_name) DO UPDATE SET current = balance.current + EXCLUDED.current;`,
		username, accrual); execErr != nil {
		b.logger.Error("Failed while making accrual", "user", username, "err", execErr)
		return execErr
	}

	return nil
}

func (b *BalanceRepository) Withdrawal(ctx context.Context, username string, withdrawal float64) error {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		b.logger.Error("Failed begin transaction while make withdrawal", "user", username, "err", err)
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
	if res, execErr = tx.ExecContext(ctx,
		`UPDATE balance SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_name = $2 AND current >= $1`,
		withdrawal, username); execErr != nil {
		b.logger.Error("Failed while making withdrawal", "user", username, "err", execErr)
		return execErr
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return balancedomain.ErrInsufficientFunds
	}

	return nil
}
