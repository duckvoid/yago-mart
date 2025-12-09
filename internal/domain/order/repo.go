package order

type Repository interface {
	All() []*Order
	Get(id int64) (*Order, error)
	GetByUser(username string) ([]*Order, error)
	Create(order *Order) error
	Update(order *Order) error
	Delete(id int64) error
}
