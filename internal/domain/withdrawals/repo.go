package withdrawals

import "context"

type Repository interface {
	GetByUser(ctx context.Context, username string) ([]*Entity, error)
	Create(ctx context.Context, withdrawal *Entity) error
}
