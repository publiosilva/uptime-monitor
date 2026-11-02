package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"uptime-monitor-backend/internal/auth"
)

type AuthMiddleware struct {
	jwtAdapter *auth.JWTAdapter
}

func NewAuthMiddleware(jwtAdapter *auth.JWTAdapter) *AuthMiddleware {
	return &AuthMiddleware{jwtAdapter: jwtAdapter}
}

func (m *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/graphql") || strings.HasSuffix(r.URL.Path, "/auth/login") || strings.HasSuffix(r.URL.Path, "/auth/register") {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userID, err := m.jwtAdapter.VerifyToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		r = r.WithContext(auth.WithUserID(r.Context(), userID))

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
