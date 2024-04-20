package geteditnodes

import (
	"context"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
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
	_, err := s.repo.EditNodesStatusAndGetReadyNodes(int(in.FreeWorkers))
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return &pb.GetNodesResponse{
		//Nodes: nodes,
	}, nil
}

func (s *Server) EditNode(ctx context.Context, in *pb.EditNodeRequest) (*empty.Empty, error) {
	return nil, nil
}
