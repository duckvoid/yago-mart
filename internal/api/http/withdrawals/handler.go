package withdrawalsapi

import (
	"net/http"

	"github.com/duckvoid/yago-mart/internal/service"
)

type WithdrawalsHandler struct {
	svc *service.UserService
}

func NewUserHandler(service *service.UserService) *WithdrawalsHandler {
	return &WithdrawalsHandler{svc: service}
}

func (u *WithdrawalsHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {}
