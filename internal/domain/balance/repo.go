package balance

type Repository interface {
	Get(username string) (*Entity, error)
	Accrual(username string, value float64) error
	Withdrawal(username string, value float64) error
}
