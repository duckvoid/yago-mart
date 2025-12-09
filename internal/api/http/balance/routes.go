package balance

import (
	"github.com/duckvoid/yago-mart/internal/api/http/middlewares"
	"github.com/go-chi/chi/v5"
)

func NewBalanceRoute(r chi.Router, handler *Handler) {
	r.Use(middlewares.AuthenticateMiddleware)

	r.Get("/balance", handler.Balance)
	r.Post("/balance/withdraw", handler.BalanceWithdraw)
}
