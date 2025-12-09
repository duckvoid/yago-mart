package withdrawals

import (
	"github.com/go-chi/chi/v5"
)

func NewWithdrawalsRoute(r chi.Router, handler *Handler) {

	r.Get("/withdrawals", handler.Withdrawals)

}
