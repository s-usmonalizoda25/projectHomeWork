package middleware

import "net/http"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized: Неверный или отсутствующий токен"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
