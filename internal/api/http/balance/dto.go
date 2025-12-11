package balance

type CurrentBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}
