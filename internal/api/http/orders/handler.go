package ordersapi

import (
	"net/http"

	"github.com/duckvoid/yago-mart/internal/service"
)

type OrdersHandler struct {
	svc *service.UserService
}

func NewOrdersHandler(service *service.UserService) *OrdersHandler {
	return &OrdersHandler{svc: service}
}

func (o *OrdersHandler) Orders(w http.ResponseWriter, r *http.Request) {}
