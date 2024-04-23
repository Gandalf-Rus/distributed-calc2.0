package geteditnodes

import (
	"context"
	"fmt"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/agent"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
	grpcconversion "github.com/Gandalf-Rus/distributed-calc2.0/internal/grpc_conversion"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	pb "github.com/Gandalf-Rus/distributed-calc2.0/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

func NewServer(repo repo) *Server {
	return &Server{
		repo: repo,
	}
}

type Server struct {
	pb.NodeServiceServer // сервис из сгенерированного пакета
	repo                 repo
}

func (s *Server) GetNodes(ctx context.Context, in *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	nodes, err := s.repo.EditNodesStatusAndGetReadyNodes(in.AgentId, int(in.FreeWorkers))
	if err != nil {
		logger.Slogger.Error(err)
		return nil, errors.ErrInternalServerError
	}

	var protoNodes []*pb.Node
	var protoNode *pb.Node
	for _, node := range nodes {
		protoNode, err = grpcconversion.NodeToGrpcNode(node)
		logger.Slogger.Info(protoNode)
		if err != nil {
			return nil, err
		}
		protoNodes = append(protoNodes, protoNode)
	}

	return &pb.GetNodesResponse{
		Nodes: protoNodes,
		OpDurations: &pb.Durations{
			Add: int32(config.Cfg.OperatorsDelay.DelayForAdd),
			Sub: int32(config.Cfg.OperatorsDelay.DelayForSub),
			Mul: int32(config.Cfg.OperatorsDelay.DelayForMul),
			Div: int32(config.Cfg.OperatorsDelay.DelayForDiv),
		},
	}, nil
}

func (s *Server) TakeHeartBeat(ctx context.Context, in *pb.GetNodesRequest) (*empty.Empty, error) {
	if agent.IsAgentRegistrated(in.AgentId) {
		agent.TakeHeartBeat(in.AgentId)
	} else {
		agent.RegistrateAgent(in.AgentId)
	}
	return nil, nil
}

func (s *Server) EditNode(ctx context.Context, in *pb.EditNodeRequest) (*empty.Empty, error) {
	node := grpcconversion.GrpcNodeToNode(in.Node)
	if err := s.repo.EditNode(node); err != nil {
		return nil, err
	}

	if node.Status == expression.Error {
		if err := s.repo.SetExpressionToError(node.ExpressionId, node.Message); err != nil {
			logger.Logger.Error(fmt.Sprintf("error to edit nodes & expression: %v", err))
		}
	} else if *node.ParentNodeId == -1 {
		if err := s.repo.SetExpressionToDone(node.ExpressionId, *node.Result); err != nil {
			logger.Logger.Error(fmt.Sprintf("error to edit expression: %v", err))
		}
	} else {
		parentNode, err := s.repo.GetNode(node.ExpressionId, *node.ParentNodeId)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("error to get parent node: %v", err))
		}

		child1, child2, err := s.repo.GetNodeChilldren(parentNode.ExpressionId, parentNode.Child1NodeId, parentNode.Child2NodeId)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("error to get children of parent node: %v", err))
		}

		if parentNode.Operand1 == nil && child1.Status == expression.Done {
			parentNode.Operand1 = child1.Result
		}
		if parentNode.Operand2 == nil && child2.Status == expression.Done {
			parentNode.Operand2 = child2.Result
		}

		logger.Logger.Info(fmt.Sprintf("parent node: nodeId: %d obj - %v", parentNode.NodeId, parentNode))

		if parentNode.Operand1 != nil && parentNode.Operand2 != nil && parentNode.Status == expression.Waiting {
			parentNode.Status = expression.Ready
			if err := s.repo.EditNode(parentNode); err != nil {
				logger.Logger.Error(fmt.Sprintf("error to edit parent node: %v", err))
			}
		}
	}

	return nil, nil
}
