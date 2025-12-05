package repository

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/duckvoid/yago-mart/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Repository struct {
	Users   *UsersRepository
	Orders  *OrdersRepository
	Balance *BalanceRepository
}

func NewRepository(ctx context.Context, database string) (*Repository, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	tableInits := []func() error{
		func() error { return initTable(ctx, db, UsersTable, embedInitUsersMigration) },
		func() error { return initTable(ctx, db, BalanceTable, embedInitBalanceMigration) },
		func() error { return initTable(ctx, db, OrdersTable, embedInitOrdersMigration) },
	}

	for _, init := range tableInits {
		if err := init(); err != nil {
			return nil, err
		}
	}

	repo := &Repository{
		Users:   NewUsersRepository(ctx, db),
		Orders:  NewOrdersRepository(ctx, db),
		Balance: NewBalanceRepository(ctx, db),
	}

	return repo, nil
}

func initTable(ctx context.Context, db *sqlx.DB, table string, embedMigration embed.FS) error {
	var exist bool

	if err := db.GetContext(ctx, &exist,
		`SELECT EXISTS(
						SELECT 1 FROM information_schema.tables 
						         WHERE table_schema = 'public' AND table_name = $1)`,
		table); err != nil {
		return err
	}

	if exist {
		return nil
	}

	logger.Log.Warn("Table %s doesn't exist. Creating...", slog.String("table", table))

	goose.SetBaseFS(embedMigration)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "init_migrations"); err != nil {
		return err
	}

	return nil
}
