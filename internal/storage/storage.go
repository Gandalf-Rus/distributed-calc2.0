package storage

import (
	"context"
	"fmt"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities"

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

	reqInsertUser = `
	INSERT INTO users (name, password) values ($1, $2)
	`

	reqSelectUserByName = `
	SELECT * FROM users WHERE name = $1
	`

	reqInsertToken = `
	INSERT INTO tokens (body) values ($1)
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

//docker run -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=distributedcalc --name distributedcalc postgres
