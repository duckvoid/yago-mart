package order

import "time"

type Status string

const (
	Invalid    Status = "INVALID"
	New        Status = "NEW"
	Processed  Status = "PROCESSED"
	Processing Status = "PROCESSING"
)

type Entity struct {
	ID          int       `db:"id"`
	Username    string    `db:"user_name"`
	Status      Status    `db:"status"`
	Accrual     float64   `db:"accrual"`
	CreatedDate time.Time `db:"created_date"`
}
