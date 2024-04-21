package geteditnodes

import (
	pb "github.com/Gandalf-Rus/distributed-calc2.0/proto"
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

// func (s *Server) GetNodes(ctx context.Context, in *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
// 	nodes, err := s.repo.EditNodesStatusAndGetReadyNodes(int(in.AgentId), int(in.FreeWorkers))
// 	if err != nil {
// 		logger.Slogger.Error(err)
// 		return nil, errors.ErrInternalServerError
// 	}

// 	var protoNodes []*pb.Node
// 	var protoNode *pb.Node
// 	for _, node := range nodes {
// 		protoNode, err = grpcconversion.NodeToGrpcNode(node)
// 		logger.Slogger.Info(protoNode)
// 		if err != nil {
// 			return nil, err
// 		}
// 		protoNodes = append(protoNodes, protoNode)
// 	}

// 	return &pb.GetNodesResponse{
// 		Nodes: protoNodes,
// 	}, nil
// }

// func (s *Server) TakeHeartBeat(ctx context.Context, in *pb.GetNodesRequest) (*empty.Empty, error) {
// 	if agent.IsAgent(int(in.AgentId)) {
// 		agent.TakeHeartBeat(int(in.AgentId))
// 	} else {
// 		agent.RegistrateAgent(int(in.AgentId))
// 	}
// 	return nil, nil
// }

// func (s *Server) EditNode(ctx context.Context, in *pb.EditNodeRequest) (*empty.Empty, error) {
// 	return nil, nil
// }
