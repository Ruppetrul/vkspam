package middleware

import "net/http"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	//Auth logic.
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
}
