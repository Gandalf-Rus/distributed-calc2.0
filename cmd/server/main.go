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
