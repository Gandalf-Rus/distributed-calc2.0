package grpcconvertion

import (
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
)

func NodeToGrpcNode(node expression.Node) proto.Node {
	return proto.Node{
		Id:           int32(node.Id),
		ExpressionId: int32(node.ExpressionId),
	}
}
