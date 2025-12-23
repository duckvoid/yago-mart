package order

import "context"

type Repository interface {
	All(ctx context.Context) ([]*Entity, error)
	Get(ctx context.Context, id int) (*Entity, error)
	GetByUser(ctx context.Context, username string) ([]*Entity, error)
	Create(ctx context.Context, order *Entity) error
	UpdateStatus(ctx context.Context, orderID int, status StatusOrder) error
	UpdateStatusAndAccrual(ctx context.Context, orderID int, accrual float64, status StatusOrder) error
}

type AccrualClient interface {
	GetOrder(ctx context.Context, orderID string) (*Accrual, error)
}
