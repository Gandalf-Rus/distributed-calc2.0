package editopduration

import (
	"errors"
	"net/http"
	"slices"

	orchErrors "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/jsonUtils"
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

		// В этом запросе достаточно проверить был ли этот токен записан в бд
		if !slices.Contains(tokens, tokenString) {
			if err = jsonUtils.RespondWith401(w, "forged or unexist JWT-token"); err != nil {
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

		err = s.Do(request)
		if err != nil {
			if errors.Is(err, orchErrors.ErrChangeConfigOperation) {
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

		respondErr := jsonUtils.SuccessRespondWith200(w, struct{}{})
		if respondErr != nil {
			logger.Slogger.Error(respondErr)
		}
	}
}
