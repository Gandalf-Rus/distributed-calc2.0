package errors

import "errors"

var (
	ErrUnuniqeUser          = errors.New("this user exists")
	ErrUnsupportedOperation = errors.New("unsupported operation")
)
