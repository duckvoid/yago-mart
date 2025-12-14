package withdrawals

import "time"

type Entity struct {
	ID          int       `db:"id"`
	Username    string    `db:"user_name"`
	OrderID     int       `db:"order_id"`
	Sum         float64   `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}
