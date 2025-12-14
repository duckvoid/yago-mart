package orders

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	orderdomain "github.com/duckvoid/yago-mart/internal/domain/order"
	"github.com/duckvoid/yago-mart/internal/service"
	"github.com/go-chi/chi/v5"
)

const maxBodySizeMib = 25

type Handler struct {
	svc *service.OrderService
}

func NewOrdersHandler(service *service.OrderService) *Handler {
	return &Handler{svc: service}
}

func (o *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySizeMib<<20)

	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)

	bodyBytes, err := io.ReadAll(tee)
	if err != nil {
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	orderID, err := strconv.Atoi(string(bodyBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !o.svc.LuhnValidation(orderID) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user := r.Context().Value("user").(string)

	if err := o.svc.Create(user, orderID); err != nil {
		switch {
		case errors.Is(err, orderdomain.ErrCreatedByAnotherUser):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, orderdomain.ErrAlreadyExist):
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(CreateResponse{
		Message: fmt.Sprintf("Order %d succesfully created", orderID),
		Code:    http.StatusAccepted,
	})
}

func (o *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)

	orders, err := o.svc.UserOrders(user)
	if err != nil {
		switch {
		case errors.Is(err, orderdomain.ErrNotFound):
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var resp ListResponse
	for _, order := range orders {
		resp.Orders = append(resp.Orders, OrderResponse{
			Number:  strconv.Itoa(order.ID),
			Status:  string(order.Status),
			Accrual: order.Accrual,
		})
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (o *Handler) Get(w http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "number")

	orderID, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !o.svc.LuhnValidation(orderID) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order, err := o.svc.Get(orderID)
	if err != nil {
		switch {
		case errors.Is(err, orderdomain.ErrNotFound):
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(OrderResponse{
		Number:  strconv.Itoa(order.ID),
		Status:  string(order.Status),
		Accrual: order.Accrual,
	})
}
