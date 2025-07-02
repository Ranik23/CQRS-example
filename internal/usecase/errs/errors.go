package errs

import "errors"



var (
	ErrNoItemsInOrder = errors.New("no items in order") //nolint:revive // this is a domain error, not a system error
	ErrInvalidUserID = errors.New("invalid user ID") //nolint:revive // this is a domain error, not a system error
	ErrInvalidOrderID = errors.New("invalid order ID") //nolint:revive // this is a domain error, not a system error
)