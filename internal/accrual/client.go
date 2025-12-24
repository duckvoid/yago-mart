package accrual

import (
	"context"
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

func (r *RestyClient) Get(ctx context.Context, url string) (*resty.Response, int, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, getTimeout)
	defer cancel()

	resp, err := r.Client.R().WithContext(timeoutCtx).Get(url)
	if err != nil {
		return nil, 0, err
	}

	return resp, resp.StatusCode(), nil
}
