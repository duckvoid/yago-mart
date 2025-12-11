package order

type Repository interface {
	All() ([]*Entity, error)
	Get(id int) (*Entity, error)
	GetByUser(username string) ([]*Entity, error)
	Create(order *Entity) error
}
