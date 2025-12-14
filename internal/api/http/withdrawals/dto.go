package withdrawals

type WithdrawalResponse struct {
	OrderID     string  `json:"order_id"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
