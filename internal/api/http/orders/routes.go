package orders

import (
	"time"

	"github.com/duckvoid/yago-mart/internal/api/http/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func NewOrdersRoute(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthenticateMiddleware)

		r.Post("/orders", handler.Create)
		r.Get("/orders", handler.List)
	})

	r.With(httprate.Limit(10, time.Minute)).Get("/order/{number}", handler.Get)

}
