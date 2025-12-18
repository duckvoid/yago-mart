package repository

import (
	"context"
	"embed"

	withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"
	"github.com/jmoiron/sqlx"
)

const WithdrawalsTable = "withdrawals"

//go:embed init_migrations/*init_withdrawals_table.sql
var embedInitWithdrawalsMigration embed.FS

type WithdrawalsRepository struct {
	db *sqlx.DB
}

func NewWithdrawalsRepository(db *sqlx.DB) *WithdrawalsRepository {
	return &WithdrawalsRepository{db: db}
}

func (w *WithdrawalsRepository) GetByUser(ctx context.Context, username string) ([]*withdrawalsdomain.Entity, error) {
	rows, err := w.db.QueryxContext(ctx, `SELECT * FROM withdrawals WHERE user_name = $1 ORDER BY processed_at`, username)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var withdrawals []*withdrawalsdomain.Entity
	for rows.Next() {
		var withdrawal withdrawalsdomain.Entity
		err = rows.StructScan(&withdrawal)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, &withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, withdrawalsdomain.ErrNotFound
	}

	return withdrawals, nil
}
