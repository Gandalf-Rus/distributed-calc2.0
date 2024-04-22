package getexpressions

import (
	"net/http"
	"slices"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/jsonUtils"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/jwt"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
)

func MakeHandler(s *Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			if err := jsonUtils.RespondWith401(w, "Authorization header is empty"); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}
		tokenString = tokenString[len("Bearer "):]
		tokens, err := s.repo.GetTokens()
		if err != nil {
			if err = jsonUtils.RespondWith500(w); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}
		if !slices.Contains(tokens, tokenString) {
			if err = jsonUtils.RespondWith401(w, "forged or unexist JWT-token"); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}

		userId, err := jwt.CheckTokenAndGetUserID(tokenString)
		if err != nil {
			if err := jsonUtils.RespondWith401(w, err.Error()); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}
		request := new(request)
		request.UserId = userId
		response, err := s.Do(request)
		if err != nil {
			logger.Slogger.Error(err)
			if err = jsonUtils.RespondWith500(w); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}

		respondErr := jsonUtils.SuccessRespondWith200(w, response)
		if respondErr != nil {
			logger.Slogger.Error(respondErr)
		}
	}
}
