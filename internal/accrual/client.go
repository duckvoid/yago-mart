package accrual

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"resty.dev/v3"
)

const getTimeout = time.Second * 5

type RestyClient struct {
	Client *resty.Client
}

func NewRestyClient() *RestyClient {
	rc := resty.New()
	return &RestyClient{Client: rc}
}

func (r *RestyClient) Get(ctx context.Context, url string) (*resty.Response, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, getTimeout)
	defer cancel()

	resp, err := r.Client.R().WithContext(timeoutCtx).Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code not 200: %d", resp.StatusCode())
	}

	return resp, nil
}
