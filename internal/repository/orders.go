package repository

import (
	"github.com/duckvoid/yago-mart/internal/model"
	"github.com/jmoiron/sqlx"
)

type OrdersRepository struct {
	db *sqlx.DB
}

func NewOrdersRepository(db *sqlx.DB) *OrdersRepository {
	return &OrdersRepository{db: db}
}

func (o *OrdersRepository) All() []*model.Order {}

func (o *OrdersRepository) Get(id int) (*model.Order, error) {}

func (o *OrdersRepository) Create(order *model.Order) error {}

func (o *OrdersRepository) Update(order *model.Order) error {}

func (o *OrdersRepository) Delete(id int) error {}
