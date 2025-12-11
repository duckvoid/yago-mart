package balance

import (
	"encoding/json"
	"errors"
	"net/http"

	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
	"github.com/duckvoid/yago-mart/internal/domain/order"
	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc *service.BalanceService
}

func NewBalanceHandler(service *service.BalanceService) *Handler {
	return &Handler{svc: service}
}

func (b *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)

	balance, err := b.svc.Get(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(CurrentBalanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	})
}

func (b *Handler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	var req WithdrawalRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(string)

	if err := b.svc.Withdrawal(user, req.OrderID, req.Sum); err != nil {
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
