package postexpression

import (
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
)

type repo interface {
	GetExpressionExitIds() ([]string, error)
	SaveExpressionAndNodes(expr expression.Expression, nodes []*expression.Node) error
}
