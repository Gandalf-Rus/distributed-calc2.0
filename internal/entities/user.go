package entities

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int64
	Name     string
	Password string
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

func (u User) ComparePasswords(inputPas string) error {
	incoming := []byte(inputPas)
	existing := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
