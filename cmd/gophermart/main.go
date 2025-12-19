package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	httpapi "github.com/duckvoid/yago-mart/internal/api/http"
	authapi "github.com/duckvoid/yago-mart/internal/api/http/auth"
	balanceapi "github.com/duckvoid/yago-mart/internal/api/http/balance"
	ordersapi "github.com/duckvoid/yago-mart/internal/api/http/orders"
	withdrawalsapi "github.com/duckvoid/yago-mart/internal/api/http/withdrawals"
	"github.com/duckvoid/yago-mart/internal/config"
	"github.com/duckvoid/yago-mart/internal/logger"
	"github.com/duckvoid/yago-mart/internal/repository"
	"github.com/duckvoid/yago-mart/internal/server"
	"github.com/duckvoid/yago-mart/internal/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}

}

func run(ctx context.Context) error {

	cfg, err := config.LoadServerConfig()
	if err != nil {
		return fmt.Errorf("failed to load server config: %w", err)
	}

	slogger := logger.New(cfg.LogLevel)

	repo, err := repository.NewRepository(ctx, cfg.Database, slogger)
	if err != nil {
		return fmt.Errorf("failed to init repository: %w", err)
	}

	userSvc := service.NewUserService(repo.Users, slogger)
	authSvc := service.NewAuthService(cfg.Secret, userSvc, slogger)
	orderSvc := service.NewOrderService(repo.Orders, slogger)
	balanceSvc := service.NewBalanceService(repo.Balance, orderSvc, slogger)
	withdrawalsSvc := service.NewWithdrawalsService(repo.Withdrawals, slogger)

	handlers := httpapi.Handlers{
		Orders:      ordersapi.NewOrdersHandler(orderSvc, slogger),
		Auth:        authapi.NewAuthHandler(authSvc, slogger),
		Balance:     balanceapi.NewBalanceHandler(balanceSvc, slogger),
		Withdrawals: withdrawalsapi.NewWithdrawalsHandler(withdrawalsSvc, slogger),
	}

	srv := server.New(cfg, handlers, slogger)

	return srv.Run(ctx)
}
