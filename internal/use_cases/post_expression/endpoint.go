package postexpression

import (
	"errors"
	"net/http"

	orchErrors "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
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
		userId, err := jwt.CheckTokenAndGetUserID(tokenString)
		if err != nil {
			if err := jsonUtils.RespondWith401(w, "Authorization header is empty"); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}

		request := new(request)
		err = jsonUtils.DecodeBody(w, r, request)
		if err != nil {
			logger.Slogger.Error("failed to decode body:", err)
			if err = jsonUtils.RespondWith400(w, "failed to decode body"); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}
		defer r.Body.Close()

		request.userId = userId
		err = s.Do(request)
		if err != nil {

			if errors.Is(err, orchErrors.ErrExistingExpression) || errors.Is(err, orchErrors.ErrIncorrectExpression) {
				if err = jsonUtils.RespondWith400(w, err.Error()); err != nil {
					logger.Slogger.Error(err)
				}
				return

			} else {
				logger.Slogger.Error(err)
				if err = jsonUtils.RespondWith500(w); err != nil {
					logger.Slogger.Error(err)
				}
				return
			}
		}

		respondErr := jsonUtils.SuccessRespondWith201(w, struct{}{})
		if respondErr != nil {
			logger.Slogger.Error(respondErr)
		}
	}
}
