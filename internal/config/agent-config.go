package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultConnHost         = "localhost"
	defaultConnPort         = 5000
	defaultMaxWorkers       = 3
	defaultHeardBeatTimeout = time.Duration(3)
	defaultGetNodesTimeout  = time.Duration(3)
)

type AgentConfig struct {
	ConnHost         string
	ConnPort         int
	MaxWorkers       int
	HeardBeatTimeout time.Duration
	GetNodesTimeout  time.Duration
	OperatorsDelay   OperatorsDelay
}

type OperatorsDelay struct {
	DelayForAdd int
	DelayForSub int
	DelayForMul int
	DelayForDiv int
}

var AgentCfg AgentConfig

func InitAgentConfig() error {
	var err error

	connHost := os.Getenv("CONN_HOST")
	if connHost == "" {
		connHost = defaultConnHost
	}

	connPort := os.Getenv("CONN_PORT")
	var port int
	if connPort == "" {
		port = defaultConnPort
	} else {
		port, err = strconv.Atoi(connPort)
		if err != nil {
			return fmt.Errorf("failed to parse %s as int: %w", os.Getenv("CONN_PORT"), err)
		}
	}

	workersCount := os.Getenv("MAX_WORKERS")
	var workers int
	if workersCount == "" {
		workers = defaultMaxWorkers
	} else {
		workers, err = strconv.Atoi(connPort)
		if err != nil {
			return fmt.Errorf("failed to parse %s as int: %w", os.Getenv("CONN_PORT"), err)
		}
	}

	AgentCfg.ConnHost = connHost
	AgentCfg.ConnPort = port
	AgentCfg.MaxWorkers = workers
	AgentCfg.HeardBeatTimeout = defaultHeardBeatTimeout
	AgentCfg.GetNodesTimeout = defaultGetNodesTimeout

	return nil
}
