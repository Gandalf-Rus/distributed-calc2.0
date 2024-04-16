package errors

import "errors"

var (
	ErrUnuniqeUser          = errors.New("this user exists")
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrIncorrectName        = errors.New("this name wasn't registrated")

	ErrExistingExpression  = errors.New("expression with this exit id already exist")
	ErrIncorrectExpression = errors.New("incorrect expression")
)
