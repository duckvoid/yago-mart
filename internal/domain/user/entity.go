package user

type Entity struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
