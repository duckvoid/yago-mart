package balanceapi

import (
	"github.com/go-chi/chi/v5"
)

func NewBalanceRoute(r chi.Router, handler *BalanceHandler) {
	r.Get("/balance", handler.Balance)
	r.Post("/balance/withdraw", handler.BalanceWithdraw)
}
