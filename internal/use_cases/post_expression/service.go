package postexpression

import (
	"slices"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/entities/expression"
	orchErr "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
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

	exitIds, err := s.repo.GetExpressionExitIds()
	if err != nil {
		return err
	}

	if slices.Contains(exitIds, request.ExitId) {
		return orchErr.ErrExistingExpression
	}

	expr, nodes, err := expression.NewExpression(request.Expression, request.ExitId)
	if err != nil {
		return orchErr.ErrIncorrectExpression
	}
	expr.UserId = request.userId

	err = s.repo.SaveExpressionAndNodes(expr, nodes)
	return err
}
