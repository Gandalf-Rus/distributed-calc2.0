package expression

import (
	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"go.uber.org/zap"
)

type Node struct {
	Id           int
	NodeId       int
	ExpressionId int
	ParentNodeId *int
	Child1NodeId *int
	Child2NodeId *int
	Operand1     *float64
	Operand2     *float64
	Operator     string
	Result       *float64
	Status       Status // (parsing, "error", waiting - ждем результатов других выражений, ready - оба операнда вычислены, in progress - передано в расчет, done - есть результат)
	Message      string
	AgentId      *string
}

type nodeStatusInfo struct { // Вспомогательная структура
	Agent_id string
	Result   float64
	Message  string
}

func (t *Expression) SetNodeStatus(node_id int, status Status, info nodeStatusInfo) {
	if int(node_id) > len(t.treeSlice)-1 || node_id < 0 {
		l.Logger.Error("Node id out of bounds",
			zap.Int("task_id", int(t.Id)),
			zap.Int("task_id", int(node_id)),
		)
		return
	}
	n := t.treeSlice[node_id]
	switch status {
	default:
		l.Logger.Error("Invalid status",
			zap.Int("task_id", int(t.Id)),
			zap.Int("node_id", int(node_id)),
			zap.String("status", status.ToString()),
		)
	case InProgress: // передано в расчет
		n.AgentId = &info.Agent_id
	case Done, Error: // есть результат или ошибка
		n.Result = &info.Result
		n.Message = info.Message
	case Parsing, Waiting, Ready:

	}
	n.Status = status
	l.Logger.Info("Node new status",
		zap.Int("task_id", t.Id),
		zap.Int("node_id", node_id),
		zap.String("status", status.ToString()),
	)

	////t.SaveTask()

	// доп. обработка после сохранения в БД
	if t.Status != InProgress {
		// делать нечего, можно забыть про результат
		return
	}
}
