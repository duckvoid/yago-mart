package user

import "context"

type Repository interface {
	All(ctx context.Context) ([]*Entity, error)
	Get(ctx context.Context, username string) (*Entity, error)
	Create(ctx context.Context, user *Entity) error
}
