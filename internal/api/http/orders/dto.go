package ordersapi

type CreateRequest struct {
	Username string `json:"username"`
	OrderID  int64  `json:"order_id"`
}

type CreateResponse struct {
	OrderID int64  `json:"order_id"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ListRequest struct {
	Username string `json:"username"`
}

type ListResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	Number     int64  `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at,omitempty"`
}
