package user

type Repository interface {
	All() ([]*Entity, error)
	Get(username string) (*Entity, error)
	Create(user *Entity) error
}
