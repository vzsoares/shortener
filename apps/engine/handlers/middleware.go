package handlers

import "net/http"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("Authing...")
		next.ServeHTTP(w, r)
	})
}
