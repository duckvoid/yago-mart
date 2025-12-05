package model

type Balance struct {
	UserID    int64
	Username  string
	Current   float64
	Withdrawn float64
}
