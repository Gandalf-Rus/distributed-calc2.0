package errors

import "errors"

var (
	ErrUnuniqeUser          = errors.New("this user exists")
	ErrUnsupportedOperation = errors.New("unsupported operation")
	IncorrectPassword       = errors.New("incorrect password")
	IncorrectName           = errors.New("this name wasn't registrated")
)
