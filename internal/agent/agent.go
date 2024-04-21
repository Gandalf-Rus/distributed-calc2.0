package agent

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	c "github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
// 	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
// 	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
// 	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type Config struct {
// 	server_host       string
// 	server_port       int
// 	max_workers       int
// 	heardBeat_timeout int
// }

// func NewConfig() Config {
// 	pserver_host := flag.String("host", c.Cfg.GrpcHost, "Host to get job from")
// 	pserver_port := flag.Int("port", c.Cfg.GrpcPort, "Port of the host")
// 	pmax_workers := flag.Int("workers", 3, "Maximum number of workers")
// 	heardBeat_timeout := flag.Int("heardBeat", int(c.Cfg.AgentLostTimeout)/2, "Poll interval (seconds)")
// 	flag.Parse()

// 	return Config{
// 		server_host:       *pserver_host,
// 		server_port:       *pserver_port,
// 		max_workers:       *pmax_workers,
// 		heardBeat_timeout: *heardBeat_timeout,
// 	}
// }

// type AgentProp struct {
// 	AgentId            string
// 	TotalProcs         int
// 	FreeProcsFreeProcs int
// }

// func New() AgentProp {
// 	agentId := 0
// 	c.Cfg.NextAgentId += 1

// 	agent := AgentProp{
// 		AgentId: int(agentId),
// 	}

// 	return agent
// }

// func (a *AgentProp) Run() {
// 	// for i := 0; i < config.max_workers; i++ {
// 	// 	wg.Add(1)
// 	// 	go Worker(choper, wg)
// 	// }

// 	host := "localhost"
// 	port := "5000"

// 	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
// 	// установим соединение
// 	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

// 	if err != nil {
// 		log.Println("could not connect to grpc server: ", err)
// 		os.Exit(1)
// 	}
// 	// закроем соединение, когда выйдем из функции
// 	defer conn.Close()

// 	grpcClient := proto.NewNodeServiceClient(conn)
// 	logger.Slogger.Info(int32(a.AgentId))
// 	grpcClient.TakeHeartBeat(context.Background(), &proto.GetNodesRequest{
// 		AgentId:     a.AgentId,
// 		FreeWorkers: 0,
// 	})
// }

// func calculateNode(node *expression.Node) {
// 	var result int
// 	var secondsDelay int

// 	switch {
// 	case node.Operand1 == nil || node.Operand2 == nil:
// 		node.Status = expression.Error
// 		node.Message = "Unready node"
// 		return
// 	case node.Operator == "+":
// 		result = *node.Operand1 + *node.Operand2
// 		node.Result = &result
// 		secondsDelay = c.Cfg.OperatorsDelay.DelayForAdd

// 	case node.Operator == "-":
// 		result = *node.Operand1 - *node.Operand2
// 		node.Result = &result
// 		secondsDelay = c.Cfg.OperatorsDelay.DelayForSub

// 	case node.Operator == "*":
// 		result = *node.Operand1 * *node.Operand2
// 		node.Result = &result
// 		secondsDelay = c.Cfg.OperatorsDelay.DelayForMul

// 	case node.Operator == "/":
// 		if *node.Operand2 == 0 {
// 			node.Status = expression.Error
// 			node.Message = "Division by zero"
// 			return
// 		} else {
// 			result = *node.Operand1 / *node.Operand2
// 			node.Result = &result
// 			secondsDelay = c.Cfg.OperatorsDelay.DelayForDiv
// 		}
// 	default:
// 		node.Status = expression.Error
// 		node.Message = "Incorrect operator [" + node.Operator + "]"
// 		log.Printf("Incorrect operator [%v] for operation %+v", node.Operator, node)
// 	}
// 	if node.Status != expression.Error {
// 		time.Sleep(time.Second * time.Duration(secondsDelay))
// 		node.Status = expression.Done
// 	}
// }
