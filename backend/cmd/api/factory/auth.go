package factory

import (
	"database/sql"
	"uptime-monitor-backend/internal/auth"
	"uptime-monitor-backend/internal/config"
)

func NewAuthHandler(db *sql.DB, cfg *config.Config) *auth.Handler {
	jwtAdapter := auth.NewJWTAdapter(cfg.JWTSecret)
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, jwtAdapter)
	authHandler := auth.NewHandler(authService)
	return authHandler
}
