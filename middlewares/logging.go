package middlewares

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

// LoggingMiddleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
