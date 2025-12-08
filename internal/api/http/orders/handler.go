package ordersapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/duckvoid/yago-mart/internal/service"
)

type OrdersHandler struct {
	svc *service.OrderService
}

func NewOrdersHandler(service *service.OrderService) *OrdersHandler {
	return &OrdersHandler{svc: service}
}

func (o *OrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if json.NewDecoder(r.Body).Decode(&req) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !o.svc.LuhnValidation(req.OrderID) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := o.svc.Create(req.Username, req.OrderID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(CreateResponse{
		Message: fmt.Sprintf("Order %d succesfully created", req.OrderID),
		Code:    http.StatusOK,
	})
}

func (o *OrdersHandler) List(w http.ResponseWriter, r *http.Request) {
	var req ListRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, err := o.svc.UserOrders(req.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resp ListResponse
	for _, order := range orders {
		resp.Orders = append(resp.Orders, OrderResponse{
			Number:     order.ID,
			Status:     string(order.Status),
			Accrual:    order.Accrual,
			UploadedAt: order.UploadDate.Format(time.RFC3339),
		})
	}

	_ = json.NewEncoder(w).Encode(resp)
}
