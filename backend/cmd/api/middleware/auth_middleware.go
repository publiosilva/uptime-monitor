package middleware

import (
	"context"
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
		if strings.HasSuffix(r.URL.Path, "/auth/login") || strings.HasSuffix(r.URL.Path, "/auth/register") {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := m.jwtAdapter.VerifyToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
