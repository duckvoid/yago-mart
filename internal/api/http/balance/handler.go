package balance

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/duckvoid/yago-mart/internal/api/http/middlewares"
	balancedomain "github.com/duckvoid/yago-mart/internal/domain/balance"
	"github.com/duckvoid/yago-mart/internal/domain/order"
	withdrawalsdomain "github.com/duckvoid/yago-mart/internal/domain/withdrawals"
	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	balanceSvc     *service.BalanceService
	withdrawalsSvc *service.WithdrawalsService
	logger         *slog.Logger
}

func NewBalanceHandler(service *service.BalanceService, withdrawalsSvc *service.WithdrawalsService, logger *slog.Logger) *Handler {
	return &Handler{balanceSvc: service, withdrawalsSvc: withdrawalsSvc, logger: logger.With(slog.String("handler", "balance"))}
}

func (b *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	user, ok := middlewares.UserFromCtx(r.Context())
	if !ok {
		b.logger.Error("failed get user from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	balance, err := b.balanceSvc.Get(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var respBuf bytes.Buffer
	if err := json.NewEncoder(&respBuf).Encode(CurrentBalanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}); err != nil {
		b.logger.Error("failed to write response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBuf.Bytes()); err != nil {
		b.logger.Error("failed to write response", "error", err)
	}
}

func (b *Handler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	var req WithdrawalRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		b.logger.Error("failed to decode request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(req.OrderID)
	if err != nil {
		b.logger.Error("failed to parse order id", slog.String("orderID", req.OrderID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := middlewares.UserFromCtx(r.Context())
	if !ok {
		b.logger.Error("failed get user from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	b.logger.Debug("Withdrawal", "request", req, "user", user)

	if err := b.balanceSvc.Withdrawal(r.Context(), user, req.Sum); err != nil {
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

	if err := b.withdrawalsSvc.Create(r.Context(), user, orderID, req.Sum); err != nil {
		switch {
		case errors.Is(err, withdrawalsdomain.ErrAlreadyExists):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

}
