package registrateuser

import (
	"strings"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities"
	orchErrs "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
)

type Service struct {
	repo repo
}

func NewSvc(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(request *request) error {

	user := entities.User{
		Name:     request.Name,
		Password: request.Password,
	}
	if err := user.SetHashedPassword(); err != nil {
		return err
	}

	err := s.repo.SaveUser(user)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return orchErrs.ErrUnuniqeUser
		}
		return err
	}

	return nil
}
