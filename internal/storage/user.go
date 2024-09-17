package storage

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/entities"

func (s *Storage) SaveUser(user entities.User) error {
	var req = `
	INSERT INTO users (name, password) values ($1, $2)
	`

	if _, err := s.connPool.Exec(s.ctx, req, user.Name, user.Password); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(name string) (entities.User, error) {
	var user entities.User

	var req = `
	SELECT * FROM users WHERE name = $1
	`

	row := s.connPool.QueryRow(s.ctx, req, name)
	err := row.Scan(&user.ID, &user.Name, &user.Password)

	return user, err
}

func (s *Storage) SaveToken(token entities.Token) error {
	var req = `
	INSERT INTO tokens (body) values ($1)
	`

	if _, err := s.connPool.Exec(s.ctx, req, token.Body); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetTokens() ([]string, error) {
	var tokens []string

	var req = `
	SELECT body FROM tokens
	`

	rows, err := s.connPool.Query(s.ctx, req)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens, err = rowsToSlice[string](rows)

	return tokens, err
}
