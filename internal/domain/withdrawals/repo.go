package withdrawals

type Repository interface {
	GetByUser(username string) ([]*Entity, error)
}
