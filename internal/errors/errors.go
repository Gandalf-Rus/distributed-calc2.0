package errors

import "errors"

var (
	ErrUnuniqeUser           = errors.New("this user exists")
	ErrUnsupportedOperation  = errors.New("unsupported operation")
	ErrIncorrectPassword     = errors.New("incorrect password")
	ErrIncorrectName         = errors.New("this name wasn't registrated")
	ErrAnotherUserExpression = errors.New("this expression belongs to another user")
	ErrUnexistExpression     = errors.New("this expression doesn't exists")

	ErrExistingExpression  = errors.New("expression with this exit id already exist")
	ErrIncorrectExpression = errors.New("incorrect expression")

	ErrInternalServerError = errors.New("internal server error")

	ErrChangeConfigOperation = errors.New("changing operation durations failed")
)
