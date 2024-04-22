package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	reqCreateUserTable = `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY NOT NULL, 
		name TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`

	reqCreateTokensTable = `
	CREATE TABLE IF NOT EXISTS tokens(
		body TEXT PRIMARY KEY NOT NULL
	);
	`

	reqCreateExpressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id SERIAL PRIMARY KEY NOT NULL,
		exit_id TEXT UNIQUE NOT NULL, 
		user_id INTEGER NOT NULL,
		body TEXT NOT NULL,
		result REAL,
		status TEXT,
		message TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	reqCreateNodesTable = `
	CREATE TABLE IF NOT EXISTS nodes(
		id             SERIAL PRIMARY KEY NOT NULL,
		node_id		   INTEGER NOT NULL,
		expression_id  INTEGER NOT NULL, 
		parent_node_id INTEGER,
		child1_node_id INTEGER,
		child2_node_id INTEGER,
		operand1       REAL,
		operand2       REAL,
		operator       CHAR NOT NULL,
		result 		   REAL,
		status         TEXT NOT NULL,
		message        TEXT,
		agent_id       TEXT,
		FOREIGN KEY (expression_id) REFERENCES expressions(id)
	);
	`

	reqInsertUser = `
	INSERT INTO users (name, password) values ($1, $2)
	`

	reqSelectUserByName = `
	SELECT * FROM users WHERE name = $1
	`

	reqInsertToken = `
	INSERT INTO tokens (body) values ($1)
	`

	reqSelectTokens = `
	SELECT body FROM tokens
	`

	reqSelectExitIds = `
	SELECT exit_id FROM expressions
	`

	reqInsertExpression = `
	INSERT INTO expressions(exit_id, user_id, body, status)
	values ($1, $2, $3, $4)
	RETURNING id
	`

	reqInsertNode = `
	INSERT INTO nodes(node_id, expression_id, parent_node_id, 
		child1_node_id, child2_node_id,
		operand1, operand2, operator, status)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	reqUpdAndGetNodesAgent = `
	UPDATE nodes SET status=$1, agent_id=$2
	WHERE id IN (SELECT id FROM nodes WHERE status='ready' LIMIT $3)
	RETURNING *
	`

	reqReleaseNodes = `
	UPDATE nodes SET status='ready', agent_id=NULL WHERE agent_id=$1
	`

	reqUpdateNode = `
	UPDATE nodes SET operand1=$1, operand2=$2, result=$3, status=$4, message=$5, agent_id=NULL 
	WHERE node_id=$6 AND expression_id=$7
	`

	reqSelectNode = `
	SELECT * FROM nodes WHERE  expression_id=$1 AND node_id=$2
	`

	reqSelectNodes = `
	SELECT * FROM nodes WHERE  expression_id=$1
	`

	reqUpdateExpression = `
	UPDATE expressions SET result=$1, status=$2, message=$3
	WHERE id=$4
	`
)

type Storage struct {
	ctx           context.Context
	connPool      *pgxpool.Pool
	closeConnOnce sync.Once
}

func New(ctx context.Context) (Storage, error) {
	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		config.Cfg.Dbuser,
		config.Cfg.Dbpassword,
		config.Cfg.Dbhost,
		config.Cfg.Dbport,
		config.Cfg.Dbname,
	)

	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return Storage{}, err
	}

	return Storage{
		ctx:           ctx,
		connPool:      conn,
		closeConnOnce: sync.Once{},
	}, nil
}

func (s *Storage) ClosePoolConn() {
	s.closeConnOnce.Do(s.connPool.Close)
}

func (s *Storage) CreateTablesIfNotExist() error {
	if err := createUsersTable(s.ctx, s.connPool); err != nil {
		return err
	}
	if err := createTokensTable(s.ctx, s.connPool); err != nil {
		return err
	}
	if err := createExpressionsTable(s.ctx, s.connPool); err != nil {
		return err
	}
	if err := createNodesTable(s.ctx, s.connPool); err != nil {
		return err
	}

	return nil
}

func (s *Storage) SaveUser(user entities.User) error {
	if _, err := s.connPool.Exec(s.ctx, reqInsertUser, user.Name, user.Password); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(name string) (entities.User, error) {
	var user entities.User

	row := s.connPool.QueryRow(s.ctx, reqSelectUserByName, name)
	err := row.Scan(&user.ID, &user.Name, &user.Password)

	return user, err
}

func (s *Storage) SaveToken(token entities.Token) error {
	if _, err := s.connPool.Exec(s.ctx, reqInsertToken, token.Body); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetExpressionExitIds() ([]string, error) {
	var ids []string

	rows, err := s.connPool.Query(s.ctx, reqSelectExitIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids, err = rowsToSlice[string](rows)

	return ids, err
}

func (s *Storage) GetTokens() ([]string, error) {
	var tokens []string

	rows, err := s.connPool.Query(s.ctx, reqSelectTokens)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens, err = rowsToSlice[string](rows)

	return tokens, err
}

func (s *Storage) SaveExpressionAndNodes(expr expression.Expression, nodes []*expression.Node) error {
	tx, err := s.connPool.BeginTx(s.ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	expressionId, err := insertExpression(s.ctx, tx, expr)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	for _, node := range nodes {
		node.ExpressionId = expressionId
		if err = insertNode(s.ctx, tx, node); err != nil {
			tx.Rollback(s.ctx)
			return err
		}
	}
	tx.Commit(s.ctx)

	return nil
}

// get edit nodes

func (s *Storage) EditNodesStatusAndGetReadyNodes(agentId string, count int) ([]*expression.Node, error) {
	tx, err := s.connPool.BeginTx(s.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	rows, err := tx.Query(s.ctx, reqUpdAndGetNodesAgent, expression.Status.ToString(expression.InProgress), agentId, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)

	if err != nil {
		return nodes, err
	}

	tx.Commit(s.ctx)
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
	if _, err := s.connPool.Exec(s.ctx, reqReleaseNodes, agentId); err != nil {
		return err
	}
	return nil
}

//
// internal methods
//

func createUsersTable(ctx context.Context, conn *pgxpool.Pool) error {
	if _, err := conn.Exec(ctx, reqCreateUserTable); err != nil {
		return err
	}
	return nil
}

func createTokensTable(ctx context.Context, conn *pgxpool.Pool) error {
	if _, err := conn.Exec(ctx, reqCreateTokensTable); err != nil {
		return err
	}
	return nil
}

func createExpressionsTable(ctx context.Context, conn *pgxpool.Pool) error {
	if _, err := conn.Exec(ctx, reqCreateExpressionsTable); err != nil {
		return err
	}
	return nil
}

func createNodesTable(ctx context.Context, conn *pgxpool.Pool) error {
	if _, err := conn.Exec(ctx, reqCreateNodesTable); err != nil {
		return err
	}
	return nil
}

func insertExpression(ctx context.Context, tx pgx.Tx, e expression.Expression) (int, error) {
	logger.Logger.Info(fmt.Sprint(e.ExitId, e.UserId, e.Body, e.Status.ToString()))
	var id int
	row := tx.QueryRow(ctx, reqInsertExpression, e.ExitId, e.UserId, e.Body, e.Status.ToString())
	err := row.Scan(&id)
	return id, err
}

func insertNode(ctx context.Context, tx pgx.Tx, n *expression.Node) error {
	if _, err := tx.Exec(ctx, reqInsertNode, n.NodeId, n.ExpressionId, n.ParentNodeId,
		n.Child1NodeId, n.Child2NodeId, n.Operand1, n.Operand2, n.Operator,
		n.Status.ToString()); err != nil {
		return err
	}

	return nil
}

func updateNode(ctx context.Context, conn *pgxpool.Pool, n *expression.Node) error {
	_, err := conn.Query(ctx, reqUpdateNode,
		n.Operand1, n.Operand2, n.Result,
		n.Status.ToString(), n.Message, n.NodeId, n.ExpressionId)

	if err != nil {
		return err
	}
	return nil
}

func updateExpression(ctx context.Context, conn *pgxpool.Pool, expressionId int, result float64, status string, message string) error {
	_, err := conn.Query(ctx, reqUpdateExpression,
		result, status, message, expressionId)

	if err != nil {
		return err
	}
	return nil
}

func getNode(ctx context.Context, conn *pgxpool.Pool, expressionId, nodeId int) (*expression.Node, error) {
	row := conn.QueryRow(ctx, reqSelectNode, expressionId, nodeId)
	node, err := scanNode(row)
	return node, err
}

func getNodes(ctx context.Context, conn *pgxpool.Pool, expressionId int) ([]*expression.Node, error) {
	rows, err := conn.Query(ctx, reqSelectNodes, expressionId)
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
