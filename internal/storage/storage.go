package storage

import (
	"context"
	"fmt"

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
		result INTEGER,
		status TEXT,
		message TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	reqCreateNodesTable = `
	CREATE TABLE IF NOT EXISTS nodes(
		id             INTEGER NOT NULL,
		expression_id  INTEGER NOT NULL, 
		parent_node_id INTEGER,
		child1_node_id INTEGER,
		child2_node_id INTEGER,
		operand1       INTEGER,
		operand2       INTEGER,
		operator       CHAR NOT NULL,
		result 		   INTEGER,
		status         TEXT,
		message        TEXT,
		agent_id       INTEGER,
		FOREIGN KEY (expression_id) REFERENCES expressions(id),
		PRIMARY KEY (id, expression_id)
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
	reqSelectExitIds = `
	SELECT exit_id FROM expressions
	`

	reqInsertExpression = `
	INSERT INTO expressions(exit_id, user_id, body, status)
	values ($1, $2, $3, $4)
	RETURNING id
	`

	reqInsertNode = `
	INSERT INTO nodes(id, expression_id, parent_node_id, 
		child1_node_id, child2_node_id,
		operand1, operand2, operator, status)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
)

type Storage struct {
	ctx context.Context
}

func New(ctx context.Context) Storage {
	return Storage{
		ctx: ctx,
	}
}

func (s Storage) CreateTablesIfNotExist() error {
	conn, err := connectToDB(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err = createUsersTable(s.ctx, conn); err != nil {
		return err
	}
	if err = createTokensTable(s.ctx, conn); err != nil {
		return err
	}
	if err = createExpressionsTable(s.ctx, conn); err != nil {
		return err
	}
	if err = createNodesTable(s.ctx, conn); err != nil {
		return err
	}

	return nil
}

func (s Storage) SaveUser(user entities.User) error {
	conn, err := connectToDB(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.Exec(s.ctx, reqInsertUser, user.Name, user.Password); err != nil {
		return err
	}
	return nil
}

func (s Storage) GetUser(name string) (entities.User, error) {
	var user entities.User

	conn, err := connectToDB(s.ctx)
	if err != nil {
		return user, err
	}
	defer conn.Close()
	row := conn.QueryRow(s.ctx, reqSelectUserByName, name)
	err = row.Scan(&user.ID, &user.Name, &user.Password)

	return user, err
}

func (s Storage) SaveToken(token entities.Token) error {
	conn, err := connectToDB(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.Exec(s.ctx, reqInsertToken, token.Body); err != nil {
		return err
	}
	return nil
}

func (s Storage) GetExpressionExitIds() ([]string, error) {
	var ids []string

	conn, err := connectToDB(s.ctx)
	if err != nil {
		return ids, err
	}
	defer conn.Close()
	row := conn.QueryRow(s.ctx, reqSelectExitIds)
	row.Scan(&ids)

	return ids, nil
}

func (s Storage) SaveExpressionAndNodes(expr expression.Expression, nodes []*expression.Node) error {
	conn, err := connectToDB(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(s.ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	expressionId, err := insertExpression(s.ctx, conn, expr)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	for _, node := range nodes {
		node.ExpressionId = expressionId
		if err = insertNode(s.ctx, conn, node); err != nil {
			tx.Rollback(s.ctx)
			return err
		}
		logger.Logger.Info("node socsess saved")
	}
	tx.Commit(s.ctx)

	return nil
}

func connectToDB(ctx context.Context) (*pgxpool.Pool, error) {
	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		config.Cfg.Dbuser,
		config.Cfg.Dbpassword,
		config.Cfg.Dbhost,
		config.Cfg.Dbport,
		config.Cfg.Dbname,
	)

	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

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

func insertExpression(ctx context.Context, conn *pgxpool.Pool, e expression.Expression) (int, error) {
	logger.Logger.Info(fmt.Sprint(e.ExitId, e.UserId, e.Body, e.Status.ToString()))
	var id int
	row := conn.QueryRow(ctx, reqInsertExpression, e.ExitId, e.UserId, e.Body, e.Status.ToString())
	err := row.Scan(&id)
	return id, err
}

func insertNode(ctx context.Context, conn *pgxpool.Pool, n *expression.Node) error {
	if _, err := conn.Exec(ctx, reqInsertNode, n.Id, n.ExpressionId, n.ParentNodeId,
		n.Child1NodeId, n.Child2NodeId, n.Operand1, n.Operand2, n.Operator,
		n.Status.ToString()); err != nil {
		return err
	}

	return nil
}

//docker run -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=distributedcalc --name distributedcalc postgres
