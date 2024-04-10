package middlewares

import (
	"net/http"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Info("Incoming request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}
