package registrateuser

import (
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities"
)

type repo interface {
	SaveUser(user entities.User) error
}
