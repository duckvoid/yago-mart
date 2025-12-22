package accrual

import "errors"

var (
	OrderNotRegistered = errors.New("order not registered")
	RateLimitExceeded  = errors.New("rate limit exceeded")
)
