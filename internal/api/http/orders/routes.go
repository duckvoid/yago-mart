package orders

import (
	"github.com/duckvoid/yago-mart/internal/api/http/middlewares"
	"github.com/go-chi/chi/v5"
)

func NewOrdersRoute(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthenticateMiddleware)

		r.Post("/orders", handler.Create)
		r.Get("/orders", handler.List)
	})

}
