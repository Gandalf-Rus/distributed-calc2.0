package getexpressions

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

	exprs, err := s.repo.GetUserExpressions(request.UserId)
	if err != nil {
		return resp, err
	}
	if exprs == nil {
		return resp, nil
	}

	for _, expr := range exprs {
		resp.Expressions = append(resp.Expressions, jsonExpression{
			Body:    expr.Body,
			Result:  expr.Result,
			Status:  expr.Status.ToString(),
			Message: expr.Message,
		})
	}

	return resp, nil
}
