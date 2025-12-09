package repository

import (
	"context"
	"embed"

	"github.com/duckvoid/yago-mart/internal/domain/order"
	"github.com/jmoiron/sqlx"
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

func (o *OrdersRepository) All() []*order.Order {
	return nil
}

func (o *OrdersRepository) Get(id int64) (*order.Order, error) {
	return nil, nil
}

func (o *OrdersRepository) GetByUser(username string) ([]*order.Order, error) {
	return nil, nil
}

func (o *OrdersRepository) Create(order *order.Order) error {
	return nil
}

func (o *OrdersRepository) Update(order *order.Order) error {
	return nil
}

func (o *OrdersRepository) Delete(id int64) error {
	return nil
}
