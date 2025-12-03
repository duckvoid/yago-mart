package repository

import (
	"context"

	"github.com/duckvoid/yago-mart/internal/model"
	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewUsersRepository(ctx context.Context, db *sqlx.DB) *UsersRepository {
	return &UsersRepository{ctx: ctx, db: db}
}

func (u *UsersRepository) All() []*model.User {}

func (u *UsersRepository) Get(id int64) (*model.User, error) {}

func (u *UsersRepository) Create(user *model.User) error {
	tx, err := u.db.BeginTxx(u.ctx, nil)
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

	if _, execErr = tx.ExecContext(u.ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2)`,
		user.Login, user.Password); execErr != nil {
		return execErr
	}

	return nil
}

func (u *UsersRepository) Update(user *model.User) error {}

func (u *UsersRepository) Delete(id int64) error {}
