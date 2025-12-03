package balanceapi

import (
	"net/http"

	"github.com/duckvoid/yago-mart/internal/service"
)

type BalanceHandler struct {
	svc *service.UserService
}

func NewBalanceHandler(service *service.UserService) *BalanceHandler {
	return &BalanceHandler{svc: service}
}

func (b *BalanceHandler) Balance(w http.ResponseWriter, r *http.Request) {}

func (b *BalanceHandler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {}
