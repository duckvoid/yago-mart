package balance

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
	"github.com/duckvoid/yago-mart/internal/domain/order"
	"github.com/duckvoid/yago-mart/internal/logger"
	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc *service.BalanceService
}

func NewBalanceHandler(service *service.BalanceService) *Handler {
	return &Handler{svc: service}
}

func (b *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	balance, err := b.svc.Get(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var respBuf bytes.Buffer
	if err := json.NewEncoder(w).Encode(CurrentBalanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBuf.Bytes()); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (b *Handler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	var req WithdrawalRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("user").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := b.svc.Withdrawal(r.Context(), user, req.OrderID, req.Sum); err != nil {
		switch {
		case errors.Is(err, order.ErrNotFound):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, balancedomain.ErrInsufficientFunds):
			http.Error(w, err.Error(), http.StatusPaymentRequired)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

}
