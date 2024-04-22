package getexpression

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"

type repo interface {
	GetTokens() ([]string, error)
	GetExpression(exitId string) (*expression.Expression, error)
}
