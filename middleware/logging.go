package middleware

import (
	"net/http"
	"project/logger"
)
func Logging(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("Входящий запрос: " + r.Method + " " + r.URL.Path)
		next.ServeHTTP(w, r)
	})
}