package main

import (
	"fmt"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
)

func main() {
	l.InitLogger()
	defer l.Logger.Sync()

	expr, nodes, err := expression.NewExpression("2+2-2/(2-1)", "0")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range nodes {
		fmt.Printf("%v. %v %v %v (%v, %v, %v)\n", v.Id, v.Operand1, v.Operator, v.Operand2, v.Child1_node_id, v.Child2_node_id, v.Status)
	}
	fmt.Println(expr.Body)
	fmt.Println(expr.Status)
}
