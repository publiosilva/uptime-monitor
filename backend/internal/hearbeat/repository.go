package hearbeat

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(heartbeat Heartbeat) (Heartbeat, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(heartbeat Heartbeat) (Heartbeat, error) {
	err := r.db.QueryRow(
		`INSERT INTO heartbeats (monitor_id, status_code, latency_ms, is_up, error_message)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, monitor_id, status_code, latency_ms, is_up, error_message, created_at`,
		heartbeat.MonitorID,
		heartbeat.StatusCode,
		heartbeat.LatencyMs,
		heartbeat.IsUp,
		heartbeat.ErrorMessage,
	).Scan(
		&heartbeat.ID,
		&heartbeat.MonitorID,
		&heartbeat.StatusCode,
		&heartbeat.LatencyMs,
		&heartbeat.IsUp,
		&heartbeat.ErrorMessage,
		&heartbeat.CreatedAt,
	)
	if err != nil {
		return Heartbeat{}, fmt.Errorf("create heartbeat: %w", err)
	}

	return heartbeat, nil
}
