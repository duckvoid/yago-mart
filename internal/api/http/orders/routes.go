package ordersapi

import (
	"github.com/go-chi/chi/v5"
)

func NewOrdersRoute(r chi.Router, handler *OrdersHandler) {
	r.Route("/", func(r chi.Router) {
		r.Post("/orders", handler.Orders)
		r.Get("/orders", handler.Orders)
	})
}
