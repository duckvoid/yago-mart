package repository

import (
	"context"
	"embed"
	"errors"
	"log/slog"

	withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const WithdrawalsTable = "withdrawals"

//go:embed init_migrations/*init_withdrawals_table.sql
var embedInitWithdrawalsMigration embed.FS

type WithdrawalsRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewWithdrawalsRepository(db *sqlx.DB, logger *slog.Logger) *WithdrawalsRepository {
	return &WithdrawalsRepository{db: db, logger: logger}
}

func (w *WithdrawalsRepository) Create(ctx context.Context, withdrawal *withdrawalsdomain.Entity) error {
	tx, err := w.db.BeginTxx(ctx, nil)
	if err != nil {
		w.logger.Error("Failed while beginning create withdrawal transaction", "error", err)
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
		`INSERT INTO withdrawals (user_name, order_id, sum) VALUES ($1, $2, $3)`,
		withdrawal.Username, withdrawal.OrderID, withdrawal.Sum); execErr != nil {
		var pgErr *pq.Error
		if errors.As(execErr, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return withdrawalsdomain.ErrAlreadyExists
			default:
				w.logger.Error("Failed while creating withdrawal", "id", withdrawal.OrderID, "user_name", withdrawal.Username, "err", execErr)
				return execErr
			}
		}
	}

	return nil
}

func (w *WithdrawalsRepository) GetByUser(ctx context.Context, username string) ([]*withdrawalsdomain.Entity, error) {
	rows, err := w.db.QueryxContext(ctx, `SELECT * FROM withdrawals WHERE user_name = $1 ORDER BY processed_at`, username)
	if err != nil {
		w.logger.Error("Failed while querying withdrawals", "user", username, "err", err)
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var withdrawals []*withdrawalsdomain.Entity
	for rows.Next() {
		var withdrawal withdrawalsdomain.Entity
		err = rows.StructScan(&withdrawal)
		if err != nil {
			w.logger.Error("Failed while scanning withdrawals", "user", username, "err", err)
			return nil, err
		}
		withdrawals = append(withdrawals, &withdrawal)
	}

	if err := rows.Err(); err != nil {
		w.logger.Error("Failed while scanning withdrawals", "user", username, "err", err)
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, withdrawalsdomain.ErrNotFound
	}

	return withdrawals, nil
}
