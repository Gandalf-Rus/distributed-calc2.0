package geteditnodes

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"

type repo interface {
	EditNodesStatusAndGetReadyNodes(agentId int, count int) ([]*expression.Node, error)
	EditNode(node *expression.Node) error
}
