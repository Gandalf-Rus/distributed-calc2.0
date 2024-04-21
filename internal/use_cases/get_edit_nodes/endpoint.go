package geteditnodes

import (
	"context"

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
	nodes, err := s.repo.EditNodesStatusAndGetReadyNodes(int(in.AgentId), int(in.FreeWorkers))
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
	}, nil
}

func (s *Server) EditNode(ctx context.Context, in *pb.EditNodeRequest) (*empty.Empty, error) {
	return nil, nil
}
