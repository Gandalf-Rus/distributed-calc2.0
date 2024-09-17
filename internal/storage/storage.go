package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
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

func createUsersTable(ctx context.Context, conn *pgxpool.Pool) error {
	var req = `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY NOT NULL, 
		name TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`

	if _, err := conn.Exec(ctx, req); err != nil {
		return err
	}
	return nil
}

func createTokensTable(ctx context.Context, conn *pgxpool.Pool) error {
	var req = `
	CREATE TABLE IF NOT EXISTS tokens(
		body TEXT PRIMARY KEY NOT NULL
	);
	`

	if _, err := conn.Exec(ctx, req); err != nil {
		return err
	}
	return nil
}

func createExpressionsTable(ctx context.Context, conn *pgxpool.Pool) error {
	var req = `
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

	if _, err := conn.Exec(ctx, req); err != nil {
		return err
	}
	return nil
}

func createNodesTable(ctx context.Context, conn *pgxpool.Pool) error {
	var req = `
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

	if _, err := conn.Exec(ctx, req); err != nil {
		return err
	}
	return nil
}
