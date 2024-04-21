package agent

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"

type SmartNode struct {
	expression.Node
	calcFunc func(node *expression.Node)
}

func NewSmartNode(node *expression.Node, calcFunc func(node *expression.Node)) *SmartNode {
	return &SmartNode{
		Node:     *node,
		calcFunc: calcFunc,
	}
}

func (sn *SmartNode) Do() {
	sn.calcFunc(&sn.Node)
}
