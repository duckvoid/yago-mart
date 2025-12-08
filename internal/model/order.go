package model

import "time"

type OrderStatus string

const (
	OrderInvalid    OrderStatus = "INVALID"
	OrderRegistered OrderStatus = "REGISTERED"
	OrderProcessed  OrderStatus = "PROCESSED"
	OrderProcessing OrderStatus = "PROCESSING"
)

type Order struct {
	ID         int
	UserID     int64
	Username   string
	Status     OrderStatus
	Accrual    int
	UploadDate time.Time
}
