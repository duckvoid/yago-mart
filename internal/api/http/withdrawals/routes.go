package withdrawals

import (
	"github.com/duckvoid/yago-mart/internal/api/http/middlewares"
	"github.com/go-chi/chi/v5"
)

func NewWithdrawalsRoute(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthenticateMiddleware)

		r.Get("/withdrawals", handler.Withdrawals)
	})
}
