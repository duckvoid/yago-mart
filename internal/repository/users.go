package repository

import (
	"context"
	"embed"
	"errors"

	userdomain "github.com/duckvoid/yago-mart/internal/domain/user"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (u *UsersRepository) All() ([]*userdomain.User, error) {
	rows, err := u.db.QueryxContext(u.ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	var users []*userdomain.User
	for rows.Next() {
		var user *userdomain.User
		err = rows.StructScan(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UsersRepository) Get(username string, password string) (*userdomain.User, error) {
	var user userdomain.User

	row := u.db.QueryRowxContext(u.ctx, `SELECT * FROM users WHERE name = $1 AND password = $2`, username, password)

	if err := row.StructScan(&user); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, userdomain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersRepository) Create(user *userdomain.User) error {
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
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return userdomain.ErrAlreadyExist
		}
		return execErr
	}

	return nil
}
