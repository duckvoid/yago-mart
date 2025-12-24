package balance

type Entity struct {
	Username  string  `db:"user_name"`
	Current   float64 `db:"current"`
	Withdrawn float64 `db:"withdrawn"`
}
