package agent

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	c "github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	grpcconversion "github.com/Gandalf-Rus/distributed-calc2.0/internal/grpc_conversion"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	AgentId        string
	TotalWorkers   int
	FreeWorkers    int32
	CtxCancelFunc  context.CancelFunc
	ctx            context.Context
	grpcClientConn *grpc.ClientConn
	serviceClient  proto.NodeServiceClient
}

func New(ctx context.Context, ctxCancelFunc context.CancelFunc) AgentProp {
	initConfig()
	agent := AgentProp{
		AgentId:       uuid.New().String(),
		TotalWorkers:  cfg.maxWorkers,
		FreeWorkers:   int32(cfg.maxWorkers),
		ctx:           ctx,
		CtxCancelFunc: ctxCancelFunc,
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

	a.serviceClient = proto.NewNodeServiceClient(conn)

	heardBeatTick := time.NewTicker(time.Second * time.Duration(cfg.heardBeatTimeout))
	TaskTick := time.NewTicker(time.Second * time.Duration(cfg.getNodesTimeout))
	chTasks := make(chan *expression.Node)
	defer close(chTasks)

	var wg *sync.WaitGroup = &sync.WaitGroup{}
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		a.getTasks(chTasks)
		done := false
		for !done {
			select {
			case <-heardBeatTick.C:
				a.doHeardBeat()
			case <-TaskTick.C:
				a.getTasks(chTasks)
			case <-a.ctx.Done():
				heardBeatTick.Stop()
				TaskTick.Stop()
				wg.Done()
				done = true
			}
		}
	}(wg)

	go func(ch <-chan *expression.Node, agent *AgentProp, wg *sync.WaitGroup) {
		// Здесь из канала приходят ноды только если у нас были свободные воркеры =>
		// тут делать проверку на свободных агентов не надо, просто читаем из канала
		for task := range ch {
			go func(task *expression.Node, a *AgentProp) {
				calculateNode(task)
				a.sendNode(task)
				// после вычеслений осовобождаем воркера
				atomic.AddInt32(&a.FreeWorkers, 1)
			}(task, a)
		}
		wg.Done()

	}(chTasks, a, wg)

	wg.Wait()
}

func (a *AgentProp) Stop(ctx context.Context) {
	logger.Logger.Info("shutdowning agent...")
	a.grpcClientConn.Close()

}

func (a *AgentProp) getNodes(freeWorkers int32) ([]*expression.Node, error) {
	nodes := make([]*expression.Node, 0)

	response, err := a.serviceClient.GetNodes(a.ctx, &proto.GetNodesRequest{
		AgentId:     a.AgentId,
		FreeWorkers: int32(freeWorkers),
	})

	if err != nil {
		return nodes, err
	}

	// получаем время выполнения
	cfg.operatorsDelay = c.OperatorsDelay{
		DelayForAdd: int(response.OpDurations.Add),
		DelayForSub: int(response.OpDurations.Sub),
		DelayForMul: int(response.OpDurations.Mul),
		DelayForDiv: int(response.OpDurations.Div),
	}

	var node *expression.Node
	for _, protoNode := range response.Nodes {
		node = grpcconversion.GrpcNodeToNode(protoNode)
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (a *AgentProp) doHeardBeat() {
	a.serviceClient.TakeHeartBeat(a.ctx, &proto.GetNodesRequest{
		AgentId:     a.AgentId,
		FreeWorkers: 0,
	})
}

func (a *AgentProp) sendNode(node *expression.Node) {
	protoNode, err := grpcconversion.NodeToGrpcNode(node)
	if err != nil {
		logger.Slogger.Error(err)
	}
	a.serviceClient.EditNode(a.ctx, &proto.EditNodeRequest{
		AgentId: a.AgentId,
		Node:    protoNode,
	})
}

func (a *AgentProp) getTasks(chTasks chan<- *expression.Node) {
	logger.Logger.Info("Get tasks", zap.Int32("freeWorkers", a.FreeWorkers))
	if a.FreeWorkers > 0 {
		nodes, err := a.getNodes(a.FreeWorkers)
		atomic.AddInt32(&a.FreeWorkers, int32(-1*len(nodes)))
		logger.Logger.Info("Get tasks", zap.Int32("freeWorkers2", a.FreeWorkers))

		if err != nil {
			logger.Logger.Error(fmt.Sprintf("can't get nodes error: %v", err))
			return
		}

		for _, node := range nodes {
			chTasks <- node
		}
	}
}

func calculateNode(node *expression.Node) {
	var result float64
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
		node.Status = expression.Done
	}
}
