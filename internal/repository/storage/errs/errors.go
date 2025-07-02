package errs

import "errors"



var (
	ErrNoOrderFound = errors.New("no order found") //nolint:revive // this is a domain error, not a system error
	ErrNoFound = errors.New("no found") //nolint:revive // this is a domain error, not a system error
)