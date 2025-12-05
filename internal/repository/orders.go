package repository

import (
	"context"
	"embed"

	"github.com/duckvoid/yago-mart/internal/model"
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

func (o *OrdersRepository) All() []*model.Order {
	return nil
}

func (o *OrdersRepository) Get(id int64) (*model.Order, error) {
	return nil, nil
}

func (o *OrdersRepository) GetByUser(username string) ([]*model.Order, error) {
	return nil, nil
}

func (o *OrdersRepository) Create(order *model.Order) error {
	return nil
}

func (o *OrdersRepository) Update(order *model.Order) error {
	return nil
}

func (o *OrdersRepository) Delete(id int64) error {
	return nil
}
