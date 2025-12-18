package withdrawals

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"
	"github.com/duckvoid/yago-mart/internal/logger"
	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc *service.WithdrawalsService
}

func NewWithdrawalsHandler(service *service.WithdrawalsService) *Handler {
	return &Handler{svc: service}
}

func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.svc.UserWithdrawals(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, withdrawalsdomain.ErrNotFound):
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var resp []WithdrawalResponse
	for _, withdrawal := range withdrawals {
		resp = append(resp, WithdrawalResponse{
			OrderID:     strconv.Itoa(withdrawal.OrderID),
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt.Format(time.RFC3339),
		})
	}

	var respBuf bytes.Buffer
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBuf.Bytes()); err != nil {
		logger.Log.Error(err.Error())
	}
}
