package order

import "context"

type Repository interface {
	All(ctx context.Context) ([]*Entity, error)
	Get(ctx context.Context, id int) (*Entity, error)
	GetByUser(ctx context.Context, username string) ([]*Entity, error)
	Create(ctx context.Context, order *Entity) error
	Update(ctx context.Context, order *Entity) error
	UpdateStatus(ctx context.Context, orderID int, status StatusOrder) error
}

type AccrualClient interface {
	GetOrder(ctx context.Context, orderID string) (*Accrual, error)
}
