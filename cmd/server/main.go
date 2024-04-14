package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/orchestrator"
)

func main() {

	l.InitLogger()
	defer l.Logger.Sync()

	// expr := expression.NewExpression("(2+2)-2/2", "0")
	// for _, v := range expr.TreeSlice {
	// 	fmt.Printf("%v. %v %v %v (%v, %v, %v)\n", v.Node_id, v.Operand1, v.Operator, v.Operand2, v.Child1_node_id, v.Child2_node_id, v.Status)
	// }
	// fmt.Println(expr.Expr_body)

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	orch, err := orchestrator.New(serverCtx)
	if err != nil {
		l.Slogger.Error(err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		orch.GracefulStop(serverCtx, sig, serverStopCtx)
	}()

	err = orch.Run()
	if err != nil {
		l.Slogger.Errorf("server start error: %v", err)
		os.Exit(1)
	}

	<-serverCtx.Done()
}
