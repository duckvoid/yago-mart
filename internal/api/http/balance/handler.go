package balanceapi

import (
	"encoding/json"
	"net/http"

	"github.com/duckvoid/yago-mart/internal/service"
)

type BalanceHandler struct {
	svc *service.BalanceService
}

func NewBalanceHandler(service *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{svc: service}
}

func (b *BalanceHandler) Balance(w http.ResponseWriter, r *http.Request) {
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

func (b *BalanceHandler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {}
