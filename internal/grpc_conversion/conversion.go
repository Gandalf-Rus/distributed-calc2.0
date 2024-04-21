package grpcconversion

import (
	"errors"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
)

func NodeToGrpcNode(node *expression.Node) (*proto.Node, error) {
	if node.ParentNodeId == nil || node.Operand1 == nil || node.Operand2 == nil {
		return &proto.Node{}, errors.New("node should not nil-fields")
	}
	return &proto.Node{
		Id:           int32(node.NodeId),
		ExpressionId: int32(node.ExpressionId),
		ParentNodeId: int32(*node.ParentNodeId),
		Operand1:     int32(*node.Operand1),
		Operand2:     int32(*node.Operand2),
		Operator:     node.Operator,
		Result:       0,
		Status:       node.Status.ToString(),
		Message:      node.Message,
	}, nil
}
