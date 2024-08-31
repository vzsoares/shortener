package handlers

import (
	"apps/engine/tools"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc, parameterStore *tools.Ssm) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingKey := r.Header.Get("X-Api-Key")
		if incomingKey == "" {
			http.Error(w, "\"Missing X-Api-Key\"", http.StatusForbidden)
			return
		}

		a4Key := parameterStore.Get("API_KEY_A4")

		if incomingKey != a4Key {
			http.Error(w, "\"Unauthorized\"", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
