package accrual

import (
	"context"
	"encoding/json"
	"errors"
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
			return nil, OrderNotRegistered
		case http.StatusTooManyRequests:
			return nil, RateLimitExceeded
		default:
			return nil, err
		}
	}

	var result OrderAccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	status, ok := result.Status.(order.StatusAccrual)
	if !ok {
		return nil, errors.New("invalid status accrual response")
	}

	return &order.Accrual{
		OrderID: result.Order,
		Status:  result.Status,
		Accrual: result.Accrual,
	}, nil
}
