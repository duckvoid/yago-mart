package accrual

import "errors"

var (
	ErrOrderNotRegistered = errors.New("order not registered")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
)
