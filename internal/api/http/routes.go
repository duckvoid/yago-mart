package http

import (
	authapi "github.com/duckvoid/yago-mart/internal/api/http/auth"
	balanceapi "github.com/duckvoid/yago-mart/internal/api/http/balance"
	ordersapi "github.com/duckvoid/yago-mart/internal/api/http/orders"
	withdrawalsapi "github.com/duckvoid/yago-mart/internal/api/http/withdrawals"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewAPIRouter(handlers Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api/user", func(r chi.Router) {
		authapi.NewAuthRoute(r, handlers.Auth)
		balanceapi.NewBalanceRoute(r, handlers.Balance)
		ordersapi.NewOrdersRoute(r, handlers.Orders)
		withdrawalsapi.NewWithdrawalsRoute(r, handlers.Withdrawals)
	})

	return r
}
