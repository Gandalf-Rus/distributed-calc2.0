package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/agent"
	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
)

func main() {
	l.InitLogger()
	defer l.Logger.Sync()

	agentCtx, agentStopCtx := context.WithCancel(context.Background())

	a := agent.New(agentCtx, agentStopCtx)
	go func() {
		a.Run()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	a.CtxCancelFunc()
	l.Logger.Info("Exit on ctrl+C signal")

	<-agentCtx.Done()
}
