package withdrawals

import (
	"net/http"

	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc *service.UserService
}

func NewWithdrawalsHandler(service *service.UserService) *Handler {
	return &Handler{svc: service}
}

func (u *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {}
