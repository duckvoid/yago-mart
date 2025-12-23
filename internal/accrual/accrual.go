package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duckvoid/yago-mart/internal/domain/order"
)

type Accrual struct {
	baseURL string
	client  *RestyClient
}

func New(baseURL string) *Accrual {
	return &Accrual{baseURL: baseURL, client: NewRestyClient()}
}

func (a *Accrual) GetOrder(ctx context.Context, orderID string) (*order.Accrual, error) {
	resp, err := a.client.Get(ctx, fmt.Sprintf("%s/api/orders/%s", a.baseURL, orderID))
	if err != nil {
		switch resp.StatusCode() {
		case http.StatusNoContent:
			return nil, ErrOrderNotRegistered
		case http.StatusTooManyRequests:
			return nil, ErrRateLimitExceeded
		default:
			return nil, err
		}
	}

	var result OrderAccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	status, err := parseAccrualStatus(result.Status)
	if err != nil {
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
		return "", fmt.Errorf("invalid status accrual response: %s", status)

	}
}
