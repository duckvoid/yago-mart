package withdrawals

import "errors"

var (
	ErrNotFound      = errors.New("withdrawals not found")
	ErrAlreadyExists = errors.New("withdrawals for same order id already exists")
)
