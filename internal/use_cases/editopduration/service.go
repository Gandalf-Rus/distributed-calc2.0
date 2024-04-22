package editopduration

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"

type Service struct {
	repo   repo
	config config
}

func NewSvc(repo repo, config config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

func (s *Service) Do(request *request) error {
	if err := s.config.ChangeOpDuration(request.Operator, request.Delay); err != nil {
		return errors.ErrChangeConfigOperation
	}
	return nil
}
