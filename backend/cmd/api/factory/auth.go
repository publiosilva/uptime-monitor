package factory

import (
	"database/sql"
	"uptime-monitor-backend/internal/auth"
)

func NewAuthHandler(db *sql.DB) *auth.Handler {
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)
	return authHandler
}
