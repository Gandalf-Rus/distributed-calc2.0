package loginuser

import (
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities"
)

type repo interface {
	GetUser(name string) (entities.User, error)
	SaveToken(token entities.Token) error
}
