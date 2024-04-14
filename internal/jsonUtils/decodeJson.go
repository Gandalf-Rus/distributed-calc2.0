package jsonUtils

import (
	"encoding/json"
	"net/http"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
)

func DecodeBody(w http.ResponseWriter, r *http.Request, request interface{}) error {
	const notValidBodyMessage = "failed to decode request"
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Slogger.Error(notValidBodyMessage, err)
		respondErr := RespondWith400(w, notValidBodyMessage)
		if respondErr != nil {
			logger.Slogger.Error("failed to write response", err)
		}
		return err
	}
	return nil
}
