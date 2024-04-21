package expression

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"

	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"go.uber.org/zap"
)

type Expression struct {
	Id        int
	ExitId    string // внешний идентификатор для идемпотентности
	UserId    int
	Body      string
	Result    *int
	Status    Status  // (parsing, error, ready, in progress, done)
	Message   string  // текстовое сообщение с результатом/ошибкой
	treeSlice []*Node // Дерево Abstract Syntax Tree
}

type ExprStatusInfo struct { // Вспомогательная структура
	Result  *int
	Message string
}

func NewExpression(expr string, ext_id string) (Expression, []*Node, error) {
	t := Expression{Body: expr, ExitId: ext_id, treeSlice: make([]*Node, 0)}
	t.SetStatus(Parsing, ExprStatusInfo{})

	root := &Node{}
	t.add(nil, root)

	parsedtree, err := parser.ParseExpr(expr)
	if err != nil {
		t.SetStatus(Error, ExprStatusInfo{Message: err.Error()})
		return t, t.treeSlice, err
	}

	err = t.buildtree(parsedtree, t.treeSlice[0])
	if err != nil {
		t.SetStatus(Error, ExprStatusInfo{Message: err.Error()})
		return t, t.treeSlice, err
	}
	t.SetStatus(Ready, ExprStatusInfo{})
	return t, t.treeSlice, err
}

func (t *Expression) buildtree(parsedtree ast.Expr, parent *Node) error {

	switch n := parsedtree.(type) {
	case *ast.BasicLit:
		//сюда попасть не должны
		l.Logger.Error("Unexpected switch case",
			zap.String("n.type", "*ast.BasicLit"))
	case *ast.BinaryExpr:
		switch n.Op {
		case token.ADD:
			parent.Operator = "+"
		case token.SUB:
			parent.Operator = "-"
		case token.MUL:
			parent.Operator = "*"
		case token.QUO:
			parent.Operator = "/"
		default:
			return unsupport(n.Op)
		}
		parent.Status = Ready // оптимистично считаем, что оба операнда будут на блечке

		switch x := n.X.(type) {
		case *ast.BasicLit:
			// вычислять не нужно
			if x.Kind != token.INT {
				return unsupport(x.Kind)
			}
			operand1, _ := strconv.Atoi(x.Value)
			parent.Operand1 = &operand1
		default:
			parent.Status = Waiting // придется вычислять операнд
			childX := t.add(&parent.NodeId, &Node{})
			errX := t.buildtree(n.X, childX)
			parent.Child1NodeId = &childX.NodeId
			if errX != nil {
				return errX
			}
		}

		switch y := n.Y.(type) {
		case *ast.BasicLit:
			// вычислять не нужно
			if y.Kind != token.INT {
				return unsupport(y.Kind)
			}
			operand2, _ := strconv.Atoi(y.Value)
			parent.Operand2 = &operand2
		default:
			parent.Status = Waiting // придется вычислять операнд
			childY := t.add(&parent.NodeId, &Node{})
			errY := t.buildtree(n.Y, childY)
			parent.Child2NodeId = &childY.NodeId
			if errY != nil {
				return errY
			}
		}
		return nil
	case *ast.ParenExpr:
		return t.buildtree(n.X, parent)
	}
	return unsupport(reflect.TypeOf(parsedtree))
}

func unsupport(i interface{}) error {
	return fmt.Errorf("%v unsupported", i)
}

func (t *Expression) add(parent_id *int, node *Node) *Node {
	node.NodeId = len(t.treeSlice)
	node.ParentNodeId = parent_id
	node.ExpressionId = t.Id
	node.Child1NodeId = nil
	node.Child2NodeId = nil
	t.treeSlice = append(t.treeSlice, node)
	return node
}

func (t *Expression) SetStatus(status Status, info ExprStatusInfo) {
	// проверим, что кто-то другой не изменил уже наш статус до нас
	if status == t.Status {
		// делать нечего
		return
	}
	switch status {
	default:
		l.Logger.Error("Invalid status",
			zap.Int("task_id", int(t.Id)),
			zap.String("status", status.ToString()),
		)
		return
	case Parsing, Ready, InProgress:
		l.Logger.Info("Task status changed",
			zap.Int("task_id", int(t.Id)),
			zap.String("status", status.ToString()),
		)
		t.Status = status
	case Done:
		t.Result = info.Result
		t.Message = fmt.Sprintf("Calculation complete. Result = %v", t.Result)
		t.Status = Done
		l.Logger.Info("Task status complete",
			zap.Int("task_id", int(t.Id)),
			zap.Int("result", *t.Result),
		)
	case Error:
		t.Message = fmt.Sprintf("Calculation failed. Error = %v", info.Message)
		t.Status = Error
		l.Logger.Error("Task failed",
			zap.Int("task_id", int(t.Id)),
			zap.String("message", info.Message),
		)
	}

	////t.SaveTask()
}

// // выбираем ожидающую операцию и переводим ее в процесс
// func (t *Expression) GetWaitingNodeAndSetProcess(agent_id string) (*Node, bool) {
// 	for _, n := range t.TreeSlice {
// 		ret := false
// 		func() {
// 			t.mx.Lock()
// 			defer t.mx.Unlock()
// 			// тут мы 100% одни
// 			if n.Status == "ready" {
// 				////n.SetToProcess(agent_id)
// 				ret = true
// 			}
// 		}()
// 		if ret {
// 			return n, true
// 		}
// 	}
// 	// нет операций готовых к вычислению
// 	return nil, false
// }
