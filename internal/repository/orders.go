package repository

import (
	"context"
	"embed"
	"errors"

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
	ctx context.Context
	db  *sqlx.DB
}

func NewOrdersRepository(ctx context.Context, db *sqlx.DB) *OrdersRepository {
	return &OrdersRepository{ctx: ctx, db: db}
}

func (o *OrdersRepository) All() ([]*orderdomain.Entity, error) {
	rows, err := o.db.QueryxContext(o.ctx, `SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	var orders []*orderdomain.Entity
	for rows.Next() {
		var order *orderdomain.Entity
		err = rows.StructScan(order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrdersRepository) Get(id int) (*orderdomain.Entity, error) {
	var order orderdomain.Entity

	row := o.db.QueryRowxContext(o.ctx, `SELECT * FROM orders WHERE id = $1`, id)

	if err := row.StructScan(&order); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orderdomain.ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (o *OrdersRepository) GetByUser(username string) ([]*orderdomain.Entity, error) {
	rows, err := o.db.QueryxContext(o.ctx, `SELECT * FROM orders WHERE user_name = $1 ORDER BY created_date`, username)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var orders []*orderdomain.Entity
	for rows.Next() {
		var order orderdomain.Entity
		err = rows.StructScan(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, orderdomain.ErrNotFound
	}

	return orders, nil
}

func (o *OrdersRepository) Create(order *orderdomain.Entity) error {
	tx, err := o.db.BeginTxx(o.ctx, nil)
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

	if _, execErr = tx.ExecContext(o.ctx,
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
				return execErr
			}
		}
	}

	return nil
}
