package balance

import "context"

type Repository interface {
	Get(ctx context.Context, username string) (*Entity, error)
	Withdrawal(ctx context.Context, username string, value float64) error
	Accrual(ctx context.Context, username string, accrual float64) error
}
