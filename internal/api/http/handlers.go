package httpapi

import (
	authapi "github.com/duckvoid/yago-mart/internal/api/http/auth"
	balanceapi "github.com/duckvoid/yago-mart/internal/api/http/balance"
	ordersapi "github.com/duckvoid/yago-mart/internal/api/http/orders"
	withdrawalsapi "github.com/duckvoid/yago-mart/internal/api/http/withdrawals"
)

type Handlers struct {
	Auth        *authapi.AuthHandler
	Balance     *balanceapi.BalanceHandler
	Orders      *ordersapi.OrdersHandler
	Withdrawals *withdrawalsapi.WithdrawalsHandler
}
