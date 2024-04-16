package expression

import (
	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"go.uber.org/zap"
)

type Node struct {
	Id             int
	Expression_id  int
	Parent_node_id int
	Child1_node_id int
	Child2_node_id int
	Operand1       Int
	Operand2       Int
	Operator       string
	OperatorDelay  int
	Status         Status // (parsing, "error", waiting - ждем результатов других выражений, ready - оба операнда вычислены, in progress - передано в расчет, done - есть результат)
	Message        string
	Agent_id       Int
	Result         int64
}

type NodeStatusInfo struct { // Вспомогательная структура
	Agent_id int64
	Result   int64
	Message  string
}

func (t *Expression) SetNodeStatus(node_id int, status Status, info NodeStatusInfo) {
	if node_id > len(t.treeSlice)-1 || node_id < 0 {
		l.Logger.Error("Node id out of bounds",
			zap.Int("task_id", t.Id),
			zap.Int("task_id", node_id),
		)
		return
	}
	n := t.treeSlice[node_id]
	switch status {
	default:
		l.Logger.Error("Invalid status",
			zap.Int("task_id", t.Id),
			zap.Int("node_id", node_id),
			zap.String("status", status.ToString()),
		)
	case InProgress: // передано в расчет
		n.Agent_id.Val = info.Agent_id
	case Done, Error: // есть результат или ошибка
		n.Result = info.Result
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

	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//

	if n.Status == Done {
		parent_id := n.Parent_node_id
		if parent_id == -1 {
			// если посчитали корневой узел, то значит выражение тоже
			t.SetStatus(Done, ExprStatusInfo{Result: info.Result})
		} else {
			parent := t.treeSlice[parent_id]
			// Запишем результат в родителя
			if parent.Child1_node_id == node_id {
				// мы - первая дочка
				parent.Operand1.IsVal = true
				parent.Operand1.Val = info.Result
			} else {
				parent.Operand2.IsVal = true
				parent.Operand2.Val = info.Result
			}
			// проверим, может и родителя можно считать?
			if parent.Status == Waiting &&
				(parent.Child1_node_id == -1 || // нет дочки
					t.treeSlice[parent.Child1_node_id].Status == Done) &&
				(parent.Child2_node_id == -1 || // нет дочки
					t.treeSlice[parent.Child2_node_id].Status == Done) {
				// дочек нет или они посчитаны, можем считать папу
				t.SetNodeStatus(parent_id, Ready, NodeStatusInfo{})
			}
		}
	} else if n.Status == Error {
		// ошибка в операци, отменяем задание и все ожидающие ноды
		t.SetStatus(Error, ExprStatusInfo{Message: info.Message})
		for _, n := range t.treeSlice {
			if n.Status == Waiting || n.Status == Ready {
				t.SetNodeStatus(n.Id, Error, NodeStatusInfo{Message: "Some other node has error"})
			}
		}
	}

}