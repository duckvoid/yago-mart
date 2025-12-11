package order

import "errors"

var (
	ErrAlreadyExist         = errors.New("order already exist")
	ErrNotFound             = errors.New("order not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrCreatedByAnotherUser = errors.New("order with same id already create by another user")
)
