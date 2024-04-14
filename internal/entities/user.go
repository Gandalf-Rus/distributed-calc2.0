package entities

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

func (u *User) SetHashedPassword() error {
	saltedBytes := []byte(u.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hash := string(hashedBytes[:])
	u.Password = hash
	return nil
}

func (u User) CheckPassword(inputPas string) error {
	incoming := []byte(inputPas)
	existing := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
