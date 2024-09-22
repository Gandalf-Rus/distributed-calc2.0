package storage

import (
	"context"
	"fmt"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *Storage) SaveLostNodes(readyStatus, progresStatus string) error {
	var req = `
	UPDATE nodes SET status=$1, agent_id=NULL
	WHERE status=$2
	`

	_, err := s.connPool.Exec(s.ctx, req, readyStatus, progresStatus)
	return err
}

func (s *Storage) GetExpressionExitIds() ([]string, error) {
	var req = `SELECT exit_id FROM expressions`
	var ids []string

	rows, err := s.connPool.Query(s.ctx, req)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids, err = rowsToSlice[string](rows)

	return ids, err
}

func (s *Storage) SaveExpressionAndNodes(expr expression.Expression, nodes []*expression.Node) error {
	// tx, err := s.connPool.BeginTx(s.ctx, pgx.TxOptions{})
	// if err != nil {
	// 	return err
	// }

	expressionId, err := insertExpression(s.ctx, s.connPool, expr)
	if err != nil {
		//tx.Rollback(s.ctx)
		return err
	}

	for _, node := range nodes {
		node.ExpressionId = expressionId
		if err = insertNode(s.ctx, s.connPool, node); err != nil {
			//tx.Rollback(s.ctx)
			return err
		}
	}
	//tx.Commit(s.ctx)

	return nil
}

// get expression

func (s *Storage) GetExpression(exitId string) (*expression.Expression, error) {
	var req = `SELECT * FROM expressions WHERE exit_id=$1`

	row := s.connPool.QueryRow(s.ctx, req, exitId)
	return scanExpression(row)
}

// get expressions

func (s *Storage) GetUserExpressions(userId int) ([]*expression.Expression, error) {
	var req = `SELECT * FROM expressions WHERE user_id=$1	`

	rows, err := s.connPool.Query(s.ctx, req, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*expression.Expression

	for rows.Next() {
		expr, err := scanExpression(rows)
		if err != nil {
			return result, err
		}
		result = append(result, expr)
	}
	return result, rows.Err()
}

// get edit nodes

func (s *Storage) EditNodesStatusAndGetReadyNodes(agentId string, count int) ([]*expression.Node, error) {
	//tx, err := s.connPool.BeginTx(s.ctx, pgx.TxOptions{})
	// if err != nil {
	// 	return nil, err
	// }
	// defer tx.Rollback(s.ctx)
	var req = `
	UPDATE nodes SET status=$1, agent_id=$2
	WHERE id IN (SELECT id FROM nodes WHERE status='ready' LIMIT $3)
	RETURNING *
	`

	rows, err := s.connPool.Query(s.ctx, req, expression.Status.ToString(expression.InProgress), agentId, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)

	if err != nil {
		return nodes, err
	}

	//tx.Commit(s.ctx)
	return nodes, nil
}

func (s *Storage) EditNode(node *expression.Node) error {
	return updateNode(s.ctx, s.connPool, node)
}

func (s *Storage) SetExpressionToError(expressionId int, message string) error {
	nodes, err := getNodes(s.ctx, s.connPool, expressionId)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		node.Status = expression.Error
		node.Message = message
		err = s.EditNode(node)
		if err != nil {
			return err
		}
	}
	if err = updateExpression(s.ctx, s.connPool, expressionId, 0, expression.Error.ToString(), message); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetNode(expressionId, nodeId int) (*expression.Node, error) {
	return getNode(s.ctx, s.connPool, expressionId, nodeId)
}

func (s *Storage) GetNodeChilldren(expressionId int, childId1, childId2 *int) (*expression.Node, *expression.Node, error) {
	var child1, child2 *expression.Node
	var err error

	if childId1 != nil {
		child1, err = getNode(s.ctx, s.connPool, expressionId, *childId1)
		if err != nil {
			return nil, nil, err
		}
	}

	if childId2 != nil {
		child2, err = getNode(s.ctx, s.connPool, expressionId, *childId2)
		if err != nil {
			return nil, nil, err
		}
	}

	return child1, child2, nil
}

func (s *Storage) SetExpressionToDone(expressionId int, result float64) error {
	nodes, err := getNodes(s.ctx, s.connPool, expressionId)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		node.Status = expression.Done
		node.Message = "super"
		err = updateNode(s.ctx, s.connPool, node)
		if err != nil {
			return err
		}
	}
	if err = updateExpression(s.ctx, s.connPool, expressionId, result, expression.Done.ToString(), "great"); err != nil {
		return err
	}

	return nil
}

// agent

func (s *Storage) ReleaseAgentUnfinishedNodes(agentId string) error {
	var req = `UPDATE nodes SET status='ready', agent_id=NULL WHERE agent_id=$1`

	if _, err := s.connPool.Exec(s.ctx, req, agentId); err != nil {
		return err
	}
	return nil
}

//
// internal methods
//

func insertExpression(ctx context.Context, conn *pgxpool.Pool, e expression.Expression) (int, error) {
	logger.Logger.Info(fmt.Sprint(e.ExitId, e.UserId, e.Body, e.Status.ToString()))
	var id int
	var req = `INSERT INTO expressions(exit_id, user_id, body, status)
		values ($1, $2, $3, $4)
		RETURNING id
	`

	row := conn.QueryRow(ctx, req, e.ExitId, e.UserId, e.Body, e.Status.ToString())
	err := row.Scan(&id)
	return id, err
}

func insertNode(ctx context.Context, conn *pgxpool.Pool, n *expression.Node) error {
	var req = `INSERT INTO nodes(node_id, expression_id, parent_node_id, 
				child1_node_id, child2_node_id,
				operand1, operand2, operator, status)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`

	if _, err := conn.Exec(ctx, req, n.NodeId, n.ExpressionId, n.ParentNodeId,
		n.Child1NodeId, n.Child2NodeId, n.Operand1, n.Operand2, n.Operator,
		n.Status.ToString()); err != nil {
		return err
	}

	return nil
}

func updateNode(ctx context.Context, conn *pgxpool.Pool, n *expression.Node) error {
	var req = `
	UPDATE nodes SET operand1=$1, operand2=$2, result=$3, status=$4, message=$5, agent_id=NULL 
	WHERE node_id=$6 AND expression_id=$7
	`

	_, err := conn.Exec(ctx, req,
		n.Operand1, n.Operand2, n.Result,
		n.Status.ToString(), n.Message, n.NodeId, n.ExpressionId)

	if err != nil {
		return err
	}
	return nil
}

func updateExpression(ctx context.Context, conn *pgxpool.Pool, expressionId int, result float64, status string, message string) error {
	var req = `
	UPDATE expressions SET result=$1, status=$2, message=$3
	WHERE id=$4
	`

	_, err := conn.Query(ctx, req,
		result, status, message, expressionId)

	if err != nil {
		return err
	}
	return nil
}

func getNode(ctx context.Context, conn *pgxpool.Pool, expressionId, nodeId int) (*expression.Node, error) {
	var req = `SELECT * FROM nodes WHERE  expression_id=$1 AND node_id=$2`

	row := conn.QueryRow(ctx, req, expressionId, nodeId)
	node, err := scanNode(row)
	return node, err
}

func getNodes(ctx context.Context, conn *pgxpool.Pool, expressionId int) ([]*expression.Node, error) {
	var req = `SELECT * FROM nodes WHERE  expression_id=$1`

	rows, err := conn.Query(ctx, req, expressionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodes(rows)
}

func scanNode(row pgx.Row) (*expression.Node, error) {
	var n expression.Node
	var status, message *string
	err := row.Scan(&n.Id, &n.NodeId,
		&n.ExpressionId, &n.ParentNodeId,
		&n.Child1NodeId, &n.Child2NodeId,
		&n.Operand1, &n.Operand2,
		&n.Operator, &n.Result,
		&status, &message, &n.AgentId)
	if err != nil {
		return nil, err
	}
	if message == nil {
		pass := ""
		message = &pass
	}
	n.Status = expression.ToStatus(*status)
	n.Message = *message
	return &n, nil
}

func scanNodes(rows pgx.Rows) ([]*expression.Node, error) {
	var nodes []*expression.Node
	for rows.Next() {
		n, err := scanNode(rows)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	if err := rows.Err(); err != nil {
		return nodes, err
	}

	return nodes, nil
}

func scanExpression(row pgx.Row) (*expression.Expression, error) {
	var expr expression.Expression
	var status string
	var msg *string
	if err := row.Scan(&expr.Id, &expr.ExitId, &expr.UserId, &expr.Body, &expr.Result, &status, &msg); err != nil {
		return nil, err
	}
	expr.Status = expression.ToStatus(status)
	if msg == nil {
		expr.Message = ""
	} else {
		expr.Message = *msg
	}

	return &expr, nil
}

func rowsToSlice[T any](rows pgx.Rows) ([]T, error) {
	var resultSlice []T
	for rows.Next() {
		var element T
		if err := rows.Scan(&element); err != nil {
			return resultSlice, err
		}
		resultSlice = append(resultSlice, element)
	}

	if err := rows.Err(); err != nil {
		return resultSlice, err
	}
	return resultSlice, nil
}

//docker run -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=distributedcalc --name distributedcalc postgres
