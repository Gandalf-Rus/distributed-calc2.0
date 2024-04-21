package agent

type repo interface {
	ReleaseAgentUnfinishedNodes(agentId string)
}
