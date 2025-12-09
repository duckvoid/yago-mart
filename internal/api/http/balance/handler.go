package balance

import (
	"encoding/json"
	"net/http"

	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc *service.BalanceService
}

func NewBalanceHandler(service *service.BalanceService) *Handler {
	return &Handler{svc: service}
}

func (b *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	balance, withdrawn, err := b.svc.Get("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(BalanceResponse{
		Current:   balance,
		Withdrawn: withdrawn,
	})
}

func (b *Handler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {}
