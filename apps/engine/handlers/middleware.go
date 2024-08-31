package handlers

import "net/http"

var A4Key string

func GetApiKey(key string) string {
    // TODO pull key
	return ""
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingKey := r.Header.Get("X-Api-Key")
		if incomingKey == "" {
			http.Error(w, "Missing X-Api-Key", http.StatusForbidden)
			return
		}

		if A4Key == "" {
			A4Key = GetApiKey("asd")
		}

		if incomingKey != A4Key {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
