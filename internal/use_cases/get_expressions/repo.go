package getexpressions

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"

type repo interface {
	GetTokens() ([]string, error)
	GetUserExpressions(userId int) ([]*expression.Expression, error)
}
