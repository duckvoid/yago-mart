package balance

import "context"

type Repository interface {
	Get(ctx context.Context, username string) (*Entity, error)
	Accrual(ctx context.Context, username string, value float64) error
	Withdrawal(ctx context.Context, username string, value float64) error
}
