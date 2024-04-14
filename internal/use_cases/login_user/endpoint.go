package loginuser

import (
	"errors"
	"net/http"

	orchErrors "github.com/Gandalf-Rus/distributed-calc2.0/internal/errors"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/jsonUtils"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
)

func MakeHandler(s *Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := new(request)
		err := jsonUtils.DecodeBody(w, r, request)
		if err != nil {
			logger.Slogger.Error("failed to decode body:", err)
			if err = jsonUtils.RespondWith400(w, "failed to decode body"); err != nil {
				logger.Slogger.Error(err)
			}
			return
		}
		defer r.Body.Close()

		response, err := s.Do(request)
		if err != nil {

			if errors.Is(err, orchErrors.IncorrectPassword) || errors.Is(err, orchErrors.IncorrectName) {
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

		respondErr := jsonUtils.SuccessRespondWith200(w, response)
		if respondErr != nil {
			logger.Slogger.Error(respondErr)
		}
	}
}
