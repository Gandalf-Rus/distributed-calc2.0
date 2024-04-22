package getexpression

import "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"

type Service struct {
	repo repo
}

func NewSvc(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(request *request) (*response, error) {
	resp := new(response)

	expr, err := s.repo.GetExpression(request.ExitId)
	if err != nil {
		return resp, err
	}
	if expr == nil {
		return resp, errors.ErrUnexistExpression
	}

	if expr.UserId == request.UserId {
		resp.Body = expr.Body
		resp.Result = expr.Result
		resp.Status = expr.Status.ToString()
		resp.Message = expr.Message
	} else {
		return resp, errors.ErrAnotherUserExpression
	}

	return resp, nil
}
