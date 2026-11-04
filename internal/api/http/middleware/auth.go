package middleware

import (
	"errors"
	"net/http"
	"strings"
)

var (
	errInvalidToken = errors.New("invalid token")
	errUnavailable  = errors.New("authentication service unavailable")
)

// Authenticate требует bearer-токен Keycloak для доступа к обработчику.
func Authenticate(verifier *Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			raw, ok := bearerToken(header)
			if !ok {
				http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			err := verifier.Verify(r.Context(), raw)
			if errors.Is(err, errUnavailable) {
				http.Error(w, "authentication service unavailable", http.StatusServiceUnavailable)
				return
			}
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// bearerToken отделяет сам токен от схемы Bearer.
func bearerToken(header string) (string, bool) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
