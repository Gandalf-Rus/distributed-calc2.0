package postexpression

import (
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
)

type repo interface {
	GetExpressionsExitIds() ([]string, error)
	SaveExpressionAndNodes(Expression expression.Expression, nodes []*expression.Node) error
}
