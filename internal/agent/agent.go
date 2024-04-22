package agent

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	c "github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	grpcconversion "github.com/Gandalf-Rus/distributed-calc2.0/internal/grpc_conversion"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/work"
	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	serverHost       string
	serverPort       int
	maxWorkers       int
	heardBeatTimeout int
	getNodesTimeout  int
	operatorsDelay   c.OperatorsDelay
}

var cfg Config

func initConfig() {
	serverHost := flag.String("host", "localhost", "Host to get job from")
	serverPort := flag.Int("port", 5000, "Port of the host")
	maxWorkers := flag.Int("workers", 3, "Maximum number of workers")
	heardBeatTimeout := flag.Int("heardBeatTimeout", 3, "heardBeat interval (seconds)")
	getNodesTimeout := flag.Int("getNodesTimeout", 3, "get tasks interval (seconds)")
	flag.Parse()

	cfg.serverHost = *serverHost
	cfg.serverPort = *serverPort
	cfg.maxWorkers = *maxWorkers
	cfg.heardBeatTimeout = *heardBeatTimeout
	cfg.getNodesTimeout = *getNodesTimeout

	cfg.operatorsDelay = c.OperatorsDelay{
		DelayForAdd: 1,
		DelayForSub: 1,
		DelayForMul: 1,
		DelayForDiv: 1,
	}
}

type AgentProp struct {
	AgentId            string
	TotalProcs         int
	FreeProcsFreeProcs int
	CtxCancelFunc      context.CancelFunc
	ctx                context.Context
	grpcClientConn     *grpc.ClientConn
}

func New(ctx context.Context, ctxCancelFunc context.CancelFunc) AgentProp {
	initConfig()
	agent := AgentProp{
		AgentId:            uuid.New().String(),
		TotalProcs:         cfg.maxWorkers,
		FreeProcsFreeProcs: cfg.maxWorkers,
		ctx:                ctx,
		CtxCancelFunc:      ctxCancelFunc,
	}

	return agent
}

func (a *AgentProp) Run() {
	addr := fmt.Sprintf("%s:%d", cfg.serverHost, cfg.serverPort)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		a.CtxCancelFunc()
	}
	a.grpcClientConn = conn

	grpcClient := proto.NewNodeServiceClient(conn)
	a.doHeardBeat(grpcClient)

	pool := work.New(cfg.maxWorkers)
	go func(pool *work.Pool) {
		pool.Run()
	}(pool)
	defer pool.Shutdown()

	go func() {
		tick := time.NewTicker(time.Second * time.Duration(cfg.getNodesTimeout))
		for range tick.C {
			tasks, err := a.getTasks(grpcClient, pool.CountOfFreeGorutines())
			if err != nil {
				logger.Slogger.Error(err)
			}
			for _, task := range tasks {
				pool.AddTask(task)
			}
		}
	}()

	<-a.ctx.Done()
}

func (a *AgentProp) Stop(ctx context.Context) {
	logger.Logger.Info("shutdowning agent...")
	a.grpcClientConn.Close()

}

func (a AgentProp) getTasks(client proto.NodeServiceClient, freeWorkers int) ([]work.Task, error) {
	tasks := make([]work.Task, 0)

	response, err := client.GetNodes(a.ctx, &proto.GetNodesRequest{
		AgentId:     a.AgentId,
		FreeWorkers: int32(freeWorkers),
	})

	if err != nil {
		return tasks, err
	}

	// получаем время выполнения
	cfg.operatorsDelay = c.OperatorsDelay{
		DelayForAdd: int(response.OpDurations.Add),
		DelayForSub: int(response.OpDurations.Sub),
		DelayForMul: int(response.OpDurations.Mul),
		DelayForDiv: int(response.OpDurations.Div),
	}

	var node *expression.Node
	doFunc := func(node *expression.Node) {
		calculateNode(node)
		logger.Logger.Info(fmt.Sprintf("node is done. %v %v %v = %v, message: %v, status: %v",
			*node.Operand1, node.Operator, *node.Operand2, *node.Result, node.Message, node.Status.ToString()))
		a.sendNode(node, client)
	}

	for _, protoNode := range response.Nodes {
		node = grpcconversion.GrpcNodeToNode(protoNode)
		logger.Logger.Info(fmt.Sprintf("node #%v%v in work (%v %v %v)", node.ExpressionId, node.NodeId, *node.Operand1, node.Operator, *node.Operand2))
		snode := NewSmartNode(node, doFunc)
		tasks = append(tasks, snode)
	}

	return tasks, nil
}

func (a AgentProp) doHeardBeat(client proto.NodeServiceClient) {
	tick := time.NewTicker(time.Second * time.Duration(cfg.heardBeatTimeout))
	go func() {
		for range tick.C {
			client.TakeHeartBeat(a.ctx, &proto.GetNodesRequest{
				AgentId:     a.AgentId,
				FreeWorkers: 0,
			})

		}
	}()
}

func (a AgentProp) sendNode(node *expression.Node, client proto.NodeServiceClient) {
	protoNode, err := grpcconversion.NodeToGrpcNode(node)
	if err != nil {
		logger.Slogger.Error(err)
	}
	client.EditNode(a.ctx, &proto.EditNodeRequest{
		AgentId: a.AgentId,
		Node:    protoNode,
	})
}

func calculateNode(node *expression.Node) {
	var result int
	var secondsDelay int

	switch {
	case node.Operand1 == nil || node.Operand2 == nil:
		node.Status = expression.Error
		node.Message = "Unready node"
		return
	case node.Operator == "+":
		result = *node.Operand1 + *node.Operand2
		node.Result = &result
		secondsDelay = cfg.operatorsDelay.DelayForAdd

	case node.Operator == "-":
		result = *node.Operand1 - *node.Operand2
		node.Result = &result
		secondsDelay = cfg.operatorsDelay.DelayForSub

	case node.Operator == "*":
		result = *node.Operand1 * *node.Operand2
		node.Result = &result
		secondsDelay = cfg.operatorsDelay.DelayForMul

	case node.Operator == "/":
		if *node.Operand2 == 0 {
			node.Status = expression.Error
			node.Message = "Division by zero"
			return
		} else {
			result = *node.Operand1 / *node.Operand2
			node.Result = &result
			secondsDelay = cfg.operatorsDelay.DelayForDiv
		}
	default:
		node.Status = expression.Error
		node.Message = "Incorrect operator [" + node.Operator + "]"
		log.Printf("Incorrect operator [%v] for operation %+v", node.Operator, node)
	}
	if node.Status != expression.Error {
		logger.Logger.Info(fmt.Sprintf("operator delay in sec: %v", secondsDelay))
		time.Sleep(time.Second * time.Duration(secondsDelay))
		logger.Logger.Info("delay end")
		node.Status = expression.Done
	}
}
