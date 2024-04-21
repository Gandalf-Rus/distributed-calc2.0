package agent

import (
	"log"
	"time"

	c "github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
)

type Agent struct {
	AgentId            int
	TotalProcs         int
	FreeProcsFreeProcs int
	LastSeen           time.Time
}

func CalculateNode(node *expression.Node) {
	var result int
	var secondsDelay int

	switch {
	case node.Operand1 == nil || node.Operand2 == nil:
		node.Status = expression.Error
		node.Message = "Unready node"
		return
	case node.Operator == "+":
		result = *node.Operand1 + *node.Operand2
		node.Result = &result
		secondsDelay = c.Cfg.OperatorsDelay.DelayForAdd

	case node.Operator == "-":
		result = *node.Operand1 - *node.Operand2
		node.Result = &result
		secondsDelay = c.Cfg.OperatorsDelay.DelayForSub

	case node.Operator == "*":
		result = *node.Operand1 * *node.Operand2
		node.Result = &result
		secondsDelay = c.Cfg.OperatorsDelay.DelayForMul

	case node.Operator == "/":
		if *node.Operand2 == 0 {
			node.Status = expression.Error
			node.Message = "Division by zero"
			return
		} else {
			result = *node.Operand1 / *node.Operand2
			node.Result = &result
			secondsDelay = c.Cfg.OperatorsDelay.DelayForDiv
		}
	default:
		node.Status = expression.Error
		node.Message = "Incorrect operator [" + node.Operator + "]"
		log.Printf("Incorrect operator [%v] for operation %+v", node.Operator, node)
	}
	if node.Status != expression.Error {
		time.Sleep(time.Second * time.Duration(secondsDelay))
		node.Status = expression.Done
	}
}
