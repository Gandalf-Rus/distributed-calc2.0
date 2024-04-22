package grpcconversion

import (
	"errors"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
)

func NodeToGrpcNode(node *expression.Node) (*proto.Node, error) {
	if node.Operand1 == nil || node.Operand2 == nil {
		return &proto.Node{}, errors.New("node should not nil-fields")
	}

	if node.ParentNodeId == nil {
		parent := -1
		node.ParentNodeId = &parent
	}

	if node.Result == nil {
		res := 0
		node.Result = &res
	}

	return &proto.Node{
		Id:           int32(node.NodeId),
		ExpressionId: int32(node.ExpressionId),
		ParentNodeId: int32(*node.ParentNodeId),
		Operand1:     int32(*node.Operand1),
		Operand2:     int32(*node.Operand2),
		Operator:     node.Operator,
		Result:       int32(*node.Result),
		Status:       node.Status.ToString(),
		Message:      node.Message,
	}, nil
}

func GrpcNodeToNode(grpcNode *proto.Node) *expression.Node {
	parent := int(grpcNode.ParentNodeId)
	operand1 := int(grpcNode.Operand1)
	operand2 := int(grpcNode.Operand2)
	result := int(grpcNode.Result)
	return &expression.Node{
		NodeId:       int(grpcNode.Id),
		ExpressionId: int(grpcNode.ExpressionId),
		ParentNodeId: &parent,
		Operand1:     &operand1,
		Operand2:     &operand2,
		Operator:     grpcNode.Operator,
		Result:       &result,
		Status:       expression.ToStatus(grpcNode.Status),
		Message:      grpcNode.Message,
	}
}
