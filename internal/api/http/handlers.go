package http

import (
	authapi "github.com/duckvoid/yago-mart/internal/api/http/auth"
	balanceapi "github.com/duckvoid/yago-mart/internal/api/http/balance"
	ordersapi "github.com/duckvoid/yago-mart/internal/api/http/orders"
	withdrawalsapi "github.com/duckvoid/yago-mart/internal/api/http/withdrawals"
)

type Handlers struct {
	Auth        *authapi.Handler
	Balance     *balanceapi.Handler
	Orders      *ordersapi.Handler
	Withdrawals *withdrawalsapi.Handler
}
