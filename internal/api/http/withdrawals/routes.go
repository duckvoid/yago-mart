package withdrawalsapi

import (
	"github.com/go-chi/chi/v5"
)

func NewWithdrawalsRoute(r chi.Router, handler *WithdrawalsHandler) {

	r.Get("/withdrawals", handler.Withdrawals)

}
