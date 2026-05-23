package middleware

import (
	"net/http"
)

type Auth struct {
	apiKey string
}

func NewAuth(key string) *Auth {
	return &Auth{apiKey: key}
}

func (a *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("X-API-Key")
		if key != a.apiKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"erro": "api key inválida ou ausente"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
