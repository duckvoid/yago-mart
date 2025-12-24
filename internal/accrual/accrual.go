package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/duckvoid/yago-mart/internal/domain/order"
)

type Accrual struct {
	baseURL string
	client  *RestyClient
	logger  *slog.Logger
}

func New(baseURL string, logger *slog.Logger) *Accrual {
	return &Accrual{baseURL: baseURL, logger: logger, client: NewRestyClient()}
}

func (a *Accrual) GetOrder(ctx context.Context, orderID string) (*order.Accrual, error) {
	resp, code, err := a.client.Get(ctx, fmt.Sprintf("%s/api/orders/%s", a.baseURL, orderID))
	if err != nil {
		a.logger.Error("Failed to get order from accrual system", "err", err)
		return nil, err
	}

	if code != http.StatusOK {
		switch code {
		case http.StatusNoContent:
			return nil, ErrOrderNotRegistered
		case http.StatusTooManyRequests:
			return nil, ErrRateLimitExceeded
		default:
			a.logger.Warn("Accrual API returned unexpected status", "code", code)
		}
	}

	var result OrderAccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		a.logger.Error("Failed to decode accrual response", "err", err)
		return nil, err
	}

	status, err := parseAccrualStatus(result.Status)
	if err != nil {
		a.logger.Error("Failed to parse accrual status", "err", err)
		return nil, err
	}

	return &order.Accrual{
		OrderID: result.Order,
		Status:  status,
		Sum:     result.Accrual,
	}, nil
}

func parseAccrualStatus(status string) (order.StatusAccrual, error) {
	switch status {
	case string(order.StatusAccrualInvalid),
		string(order.StatusAccrualRegistred),
		string(order.StatusAccrualProcessing),
		string(order.StatusAccrualProcessed):
		return order.StatusAccrual(status), nil
	default:
		return "", ErrUnexpectedStatus

	}
}
