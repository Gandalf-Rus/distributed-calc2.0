package loginuser

import (
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities"
	orchErrors "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/jwt"
)

type Service struct {
	repo repo
}

func NewSvc(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(request *request) (response, error) {

	user, err := s.repo.GetUser(request.Name)
	if err != nil {
		return response{}, orchErrors.IncorrectName
	}

	if err := user.CheckPassword(request.Password); err != nil {
		return response{}, orchErrors.IncorrectPassword
	}

	resp := response{}
	body, err := jwt.CreateToken(user.ID)
	if err != nil {
		return resp, err
	}
	resp.JwtToken = body

	token := entities.Token{
		Body: body,
	}
	if err := s.repo.SaveToken(token); err != nil {
		return resp, err
	}

	return resp, nil
}
