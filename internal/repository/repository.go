package repository

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Repository struct {
	Users       *UsersRepository
	Orders      *OrdersRepository
	Balance     *BalanceRepository
	Withdrawals *WithdrawalsRepository
}

func NewRepository(ctx context.Context, database string, logger *slog.Logger) (*Repository, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	tableInits := []func() error{
		func() error { return initTable(ctx, db, UsersTable, embedInitUsersMigration, logger) },
		func() error { return initTable(ctx, db, BalanceTable, embedInitBalanceMigration, logger) },
		func() error { return initTable(ctx, db, OrdersTable, embedInitOrdersMigration, logger) },
		func() error { return initTable(ctx, db, WithdrawalsTable, embedInitWithdrawalsMigration, logger) },
	}

	for _, init := range tableInits {
		if err := init(); err != nil {
			logger.Error("Failed to initialize table", "error", err)
			return nil, err
		}
	}

	repo := &Repository{
		Users:       NewUsersRepository(db, logger),
		Orders:      NewOrdersRepository(db, logger),
		Balance:     NewBalanceRepository(db, logger),
		Withdrawals: NewWithdrawalsRepository(db, logger),
	}

	return repo, nil
}

func initTable(ctx context.Context, db *sqlx.DB, table string, embedMigration embed.FS, logger *slog.Logger) error {
	var exist bool

	if err := db.GetContext(ctx, &exist,
		`SELECT EXISTS(
						SELECT 1 FROM information_schema.tables 
						         WHERE table_schema = 'public' AND table_name = $1)`,
		table); err != nil {
		logger.Error("Failed to check if table exists", "error", err)
		return err
	}

	if exist {
		logger.Info("Table already exists", "table", table)
		return nil
	}

	logger.Warn("Table doesn't exist. Creating...", slog.String("table", table))

	goose.SetBaseFS(embedMigration)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("Failed to set postgres dialect", "error", err)
		return err
	}

	if err := goose.Up(db.DB, "init_migrations"); err != nil {
		logger.Error("Failed to up migrations", "error", err)
		return err
	}

	return nil
}
