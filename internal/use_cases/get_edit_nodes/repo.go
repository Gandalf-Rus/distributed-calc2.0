package geteditnodes

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"

type repo interface {
	EditNodesStatusAndGetReadyNodes(agentId string, count int) ([]*expression.Node, error)
	EditNode(node *expression.Node) error
	SetExpressionToError(expressionId int, message string) error
	GetNode(expressionId, nodeId int) (*expression.Node, error)
	GetNodeChilldren(expressionId int, childId1, childId2 *int) (*expression.Node, *expression.Node, error)
	SetExpressionToDone(expressionId int, result float64) error
}
