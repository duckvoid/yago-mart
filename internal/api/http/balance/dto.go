package balanceapi

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalRequest struct{}
type WithdrawResponse struct{}
