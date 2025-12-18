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
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (u *UsersRepository) All(ctx context.Context) ([]*userdomain.Entity, error) {
	rows, err := u.db.QueryxContext(ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	var users []*userdomain.Entity
	for rows.Next() {
		var user *userdomain.Entity
		err = rows.StructScan(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UsersRepository) Get(ctx context.Context, username string) (*userdomain.Entity, error) {
	var user userdomain.Entity

	row := u.db.QueryRowxContext(ctx, `SELECT * FROM users WHERE name = $1`, username)

	if err := row.StructScan(&user); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, userdomain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersRepository) Create(ctx context.Context, user *userdomain.Entity) error {
	tx, err := u.db.BeginTxx(ctx, nil)
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

	if _, execErr = tx.ExecContext(ctx,
		`INSERT INTO users (name, password) VALUES ($1, $2)`,
		user.Name, user.Password); execErr != nil {
		var pgErr *pq.Error
		if errors.As(execErr, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return userdomain.ErrAlreadyExist
		}
		return execErr
	}

	return nil
}
