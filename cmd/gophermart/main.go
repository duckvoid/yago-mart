package main

import (
	"context"
	"fmt"
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

	cfg, err := config.LoadServerConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load server config: %w", err))
	}

	logger.New(cfg.LogLevel)

	repo, err := repository.NewRepository(ctx, cfg.Database)
	if err != nil {
		panic(fmt.Errorf("failed to init repository: %w", err))
	}

	userSvc := service.NewUserService(repo.Users)
	authSvc := service.NewAuthService(userSvc)
	orderSvc := service.NewOrderService(repo.Orders)
	balanceSvc := service.NewBalanceService(repo.Balance, orderSvc)
	withdrawalsSvc := service.NewWithdrawalsService(repo.Withdrawals)

	handlers := httpapi.Handlers{
		Orders:      ordersapi.NewOrdersHandler(orderSvc),
		Auth:        authapi.NewAuthHandler(authSvc),
		Balance:     balanceapi.NewBalanceHandler(balanceSvc),
		Withdrawals: withdrawalsapi.NewWithdrawalsHandler(withdrawalsSvc),
	}

	srv := server.New(cfg, handlers)

	srv.Run(ctx)
}
