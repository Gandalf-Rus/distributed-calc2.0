package agent

type repo interface {
	ReleaseAgentUnfinishedNodes(agentId int)
}
