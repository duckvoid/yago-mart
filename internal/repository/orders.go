package repository

import (
	"context"
	"embed"
	"errors"
	"log/slog"

	orderdomain "github.com/duckvoid/yago-mart/internal/domain/order"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const OrdersTable = "orders"

//go:embed init_migrations/*init_orders_table.sql
var embedInitOrdersMigration embed.FS

type OrdersRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewOrdersRepository(db *sqlx.DB, logger *slog.Logger) *OrdersRepository {
	return &OrdersRepository{db: db, logger: logger}
}

func (o *OrdersRepository) All(ctx context.Context) ([]*orderdomain.Entity, error) {
	rows, err := o.db.QueryxContext(ctx, `SELECT * FROM orders`)
	if err != nil {
		o.logger.Error("Failed while querying all orders", "error", err)
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	if err := rows.Err(); err != nil {
		o.logger.Error("Failed while iteration all orders rows", "error", err)
		return nil, err
	}

	var orders []*orderdomain.Entity
	for rows.Next() {
		var order *orderdomain.Entity
		err = rows.StructScan(order)
		if err != nil {
			o.logger.Error("Failed while scanning all orders struct", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrdersRepository) Get(ctx context.Context, id int) (*orderdomain.Entity, error) {
	var order orderdomain.Entity

	row := o.db.QueryRowxContext(ctx, `SELECT * FROM orders WHERE id = $1`, id)

	if err := row.StructScan(&order); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orderdomain.ErrNotFound
		}
		o.logger.Error("Failed while scanning order row", "id", id, "error", err)
		return nil, err
	}

	return &order, nil
}

func (o *OrdersRepository) GetByUser(ctx context.Context, username string) ([]*orderdomain.Entity, error) {
	rows, err := o.db.QueryxContext(ctx, `SELECT * FROM orders WHERE user_name = $1 ORDER BY created_date`, username)
	if err != nil {
		o.logger.Error("Failed while querying orders by user", "username", username, "error", err)
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var orders []*orderdomain.Entity
	for rows.Next() {
		var order orderdomain.Entity
		err = rows.StructScan(&order)
		if err != nil {
			o.logger.Error("Failed while scanning order row", "username", username, "error", err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		o.logger.Error("Failed while scanning order rows", "error", err)
		return nil, err
	}

	if len(orders) == 0 {
		return nil, orderdomain.ErrNotFound
	}

	return orders, nil
}

func (o *OrdersRepository) Create(ctx context.Context, order *orderdomain.Entity) error {
	tx, err := o.db.BeginTxx(ctx, nil)
	if err != nil {
		o.logger.Error("Failed while beginning create order transaction", "error", err)
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
		`INSERT INTO orders (id, user_name, status) VALUES ($1, $2, $3)`,
		order.ID, order.Username, order.Status); execErr != nil {
		var pgErr *pq.Error
		if errors.As(execErr, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return orderdomain.ErrAlreadyExist
			case pgerrcode.InvalidColumnReference:
				return orderdomain.ErrUserNotFound
			default:
				o.logger.Error("Failed while creating order", "id", order.ID, "user_name", order.Username)
				return execErr
			}
		}
	}

	return nil
}

func (o *OrdersRepository) UpdateStatus(ctx context.Context, orderID int, status orderdomain.StatusOrder) error {
	tx, err := o.db.BeginTxx(ctx, nil)
	if err != nil {
		o.logger.Error("Failed while beginning create order transaction", "error", err)
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

	if _, err := tx.ExecContext(ctx,
		`UPDATE orders SET status = $1 WHERE id = $2`, string(status), orderID); err != nil {
		o.logger.Error("Failed while updating order status ", "id", orderID, "status", status)
		return err
	}

	return nil
}

func (o *OrdersRepository) UpdateStatusAndAccrual(ctx context.Context, orderID int, accrual float64, status orderdomain.StatusOrder) error {
	tx, err := o.db.BeginTxx(ctx, nil)
	if err != nil {
		o.logger.Error("Failed while beginning create order transaction", "error", err)
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

	if _, err := tx.ExecContext(ctx,
		`UPDATE orders SET status = $1, accrual = $2 WHERE id = $3`, string(status), accrual, orderID); err != nil {
		o.logger.Error("Failed while updating order status and accrual", "id", orderID, "status", status)
		return err
	}

	return nil
}
