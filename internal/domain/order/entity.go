package order

import "time"

type StatusOrder string

type StatusAccrual string

const (
	StatusOrderInvalid    StatusOrder = "INVALID"
	StatusOrderNew        StatusOrder = "NEW"
	StatusOrderProcessed  StatusOrder = "PROCESSED"
	StatusOrderProcessing StatusOrder = "PROCESSING"
)

const (
	StatusAccrualRegistred  StatusAccrual = "REGISTRED"
	StatusAccrualInvalid    StatusAccrual = "INVALID"
	StatusAccrualProcessed  StatusAccrual = "PROCESSED"
	StatusAccrualProcessing StatusAccrual = "PROCESSING"
)

type Entity struct {
	ID          int         `db:"id"`
	Username    string      `db:"user_name"`
	Status      StatusOrder `db:"status"`
	Accrual     float64     `db:"accrual"`
	CreatedDate time.Time   `db:"created_date"`
}

type Accrual struct {
	OrderID string
	Status  StatusAccrual
	Sum     float64
}
