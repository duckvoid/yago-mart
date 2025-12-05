package ordersapi

import (
	"github.com/go-chi/chi/v5"
)

func NewOrdersRoute(r chi.Router, handler *OrdersHandler) {
	r.Post("/orders", handler.Create)
	r.Get("/orders", handler.List)
}
