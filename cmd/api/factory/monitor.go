package factory

import (
	"database/sql"
	"uptime-monitor-backend/internal/monitor"
)

func NewMonitorHandler(db *sql.DB) *monitor.Handler {
	monitorRepo := monitor.NewRepository(db)
	monitorService := monitor.NewService(monitorRepo)
	return monitor.NewHandler(monitorService)
}
