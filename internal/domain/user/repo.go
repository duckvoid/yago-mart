package user

type Repository interface {
	All() ([]*User, error)
	Get(username string, password string) (*User, error)
	Create(user *User) error
}
