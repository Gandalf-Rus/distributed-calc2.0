package agent

import (
	"time"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"go.uber.org/zap"
)

var Agents map[int]Agent = make(map[int]Agent)

// Удаляем пропавших агентов
func CleanLostAgents(repo repo) {
	for _, a := range Agents {
		if time.Since(a.LastSeen) > config.Cfg.AgentLostTimeout {
			logger.Logger.Info("Agent lost",
				zap.Int("agent_id", a.AgentId),
				zap.Time("Last seen", a.LastSeen),
				zap.Int("timeout sec", int(config.Cfg.AgentLostTimeout)),
			)

			repo.ReleaseAgentUnfinishedNodes(a.AgentId)
			delete(Agents, a.AgentId)
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
