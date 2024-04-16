package agent

// type Agent struct {
// 	AgentId    string `json:"agent_id"`
// 	Status     string `json:"status"`
// 	TotalProcs int    `json:"total_procs"`
// 	IdleProcs  int    `json:"idle_procs"`
// 	FirstSeen  time.Time
// 	LastSeen   time.Time
// }

// var Agents map[string]Agent = make(map[string]Agent)

// // Удаляем пропавших агентов
// func CleanLostAgents() {
// 	timeout := time.Second * time.Duration(config.Cfg.AgentLostTimeout)
// 	for _, a := range Agents {
// 		if time.Since(a.LastSeen) > timeout {
// 			// давно не видели, забудем про него
// 			l.Logger.Info("Agent lost",
// 				zap.String("agent_id", a.AgentId),
// 				zap.Time("Last seen", a.LastSeen),
// 				zap.Int("timeout sec", config.Cfg.AgentLostTimeout),
// 			)
// 			// но вначале передадим его задание другим
// 			// for _, t := range expression.Expressions {
// 			// 	if t.Status == "in progress" {
// 			// 		for _, n := range t.TreeSlice {
// 			// 			if n.Status == "in progress" &&
// 			// 				n.Agent_id == a.AgentId {
// 			// 				t.SetNodeStatus(n.Node_id, "ready", expression.NodeStatusInfo{})
// 			// 			}
// 			// 		}
// 			// 	}
// 			// }
// 			// нет больше такого агента
// 			delete(Agents, a.AgentId)
// 		}
// 	}
// }

// func CleanLoopAgents() {
// 	tick := time.NewTicker(time.Second * time.Duration(config.Cfg.AgentLostTimeout))
// 	go func() {
// 		for range tick.C {
// 			// таймер прозвенел
// 			CleanLostAgents()
// 		}
// 	}()
// }

// func CalculateNode(node *expression.Node) {
// 	switch {
// 	case !node.Operand1.IsVal || !node.Operand2.IsVal:
// 		node.Status = expression.Error
// 		node.Message = "Unready node"
// 	case node.Operator == "+":
// 		node.Result = node.Operand1.Val + node.Operand2.Val
// 		//operation.Result, no_overfl = overflow.Add64(int64(operation.Operand1), int64(operation.Operand2))
// 	case node.Operator == "-":
// 		node.Result = node.Operand1.Val - node.Operand2.Val
// 		//operation.Result, no_overfl = overflow.Sub64(int64(operation.Operand1), int64(operation.Operand2))
// 	case node.Operator == "*":
// 		node.Result = node.Operand1.Val * node.Operand2.Val
// 		//operation.Result, no_overfl = overflow.Mul64(int64(operation.Operand1), int64(operation.Operand2))
// 	case node.Operator == "/":
// 		if node.Operand2.Val == 0 {
// 			node.Status = expression.Error
// 			node.Message = "Division by zero"
// 		} else {
// 			node.Result = node.Operand1.Val / node.Operand2.Val
// 			//operation.Result, no_overfl = overflow.Div64(int64(operation.Operand1), int64(operation.Operand2))
// 		}
// 	default:
// 		node.Status = expression.Error
// 		node.Message = "Incorrect operator [" + node.Operator + "]"
// 		log.Printf("Incorrect operator [%v] for operation %+v", node.Operator, node)
// 	}
// 	// if !no_overfl {
// 	// 	node.Status = "error"
// 	// 	node.Message = "Overflow"
// 	// }
// 	if node.Status != expression.Error {
// 		// изображаем бурную деятельность
// 		time.Sleep(time.Duration(node.OperatorDelay) * time.Second)
// 		node.Status = expression.Done
// 	}
// }
