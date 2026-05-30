package middleware

import (
	"fmt"
	"net/http"
)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[ЛОГ] %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
