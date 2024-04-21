package agent

import (
	"fmt"
	"time"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"go.uber.org/zap"
)

type agent struct {
	LastSeen time.Time
}

var agents map[string]*agent = make(map[string]*agent)

func RegistrateAgent(agentId string) {
	agents[agentId] = &agent{
		LastSeen: time.Now(),
	}
	logger.Logger.Info(fmt.Sprintf("%v", agents))
}

func IsAgent(agentId string) bool {
	_, found := agents[agentId]
	return found
}

func TakeHeartBeat(agentId string) {
	agents[agentId].LastSeen = time.Now()
}

// Удаляем пропавших агентов
func CleanLostAgents(repo repo) {
	for id, a := range agents {
		if time.Since(a.LastSeen) > config.Cfg.AgentLostTimeout {
			logger.Logger.Info("Agent lost",
				zap.String("agent_id", id),
				zap.Time("Last seen", a.LastSeen),
				zap.Int("timeout sec", int(config.Cfg.AgentLostTimeout)),
			)

			repo.ReleaseAgentUnfinishedNodes(id)
			delete(agents, id)
		}
	}
}

func LostAgentCollector(repo repo) {
	tick := time.NewTicker(time.Second * time.Duration(config.Cfg.AgentLostTimeout))
	go func() {
		for range tick.C {
			// таймер прозвенел
			CleanLostAgents(repo)
		}
	}()
}
