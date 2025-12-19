package withdrawals

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/duckvoid/yago-mart/internal/api/http/middlewares"
	withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"
	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc    *service.WithdrawalsService
	logger *slog.Logger
}

func NewWithdrawalsHandler(service *service.WithdrawalsService, logger *slog.Logger) *Handler {
	return &Handler{svc: service, logger: logger.With(slog.String("handler", "withdrawals"))}
}

func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	user, ok := middlewares.UserFromCtx(r.Context())
	if !ok {
		h.logger.Error("failed to get user from context")
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
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBuf.Bytes()); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}
