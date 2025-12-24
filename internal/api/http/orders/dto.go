package orders

type CreateRequest struct {
	Username string `json:"username"`
	OrderID  int    `json:"order_id"`
}

type CreateResponse struct {
	OrderID int    `json:"order_id"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ListRequest struct {
	Username string `json:"username"`
}

type OrderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}
