package repository

import (
	"context"
	"embed"

	"github.com/duckvoid/yago-mart/internal/model"
	"github.com/jmoiron/sqlx"
)

const UsersTable = "users"

//go:embed init_migrations/*init_users_table.sql
var embedInitUsersMigration embed.FS

type UsersRepository struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewUsersRepository(ctx context.Context, db *sqlx.DB) *UsersRepository {
	return &UsersRepository{ctx: ctx, db: db}
}

func (u *UsersRepository) All() ([]*model.User, error) {
	rows, err := u.db.QueryxContext(u.ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	var users []*model.User
	for rows.Next() {
		var user *model.User
		err = rows.StructScan(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UsersRepository) Get(username string) (*model.User, error) {
	var user model.User

	row := u.db.QueryRowxContext(u.ctx, `SELECT * FROM users WHERE name = $1`, username)

	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

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
		`INSERT INTO users (name, password) VALUES ($1, $2)`,
		user.Name, user.Password); execErr != nil {
		return execErr
	}

	return nil
}
