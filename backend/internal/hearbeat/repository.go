package hearbeat

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(heartbeat Heartbeat) (Heartbeat, error)
	List24hByMonitorID(monitorID string) ([]Heartbeat, error)
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

func (r *PostgresRepository) List24hByMonitorID(monitorID string) ([]Heartbeat, error) {
	rows, err := r.db.Query(
		`SELECT id, monitor_id, status_code, latency_ms, is_up, error_message, created_at
		FROM heartbeats
		WHERE monitor_id = $1
		AND created_at >= NOW() - INTERVAL '24 hours'
		ORDER BY created_at DESC`,
		monitorID,
	)
	if err != nil {
		return nil, fmt.Errorf("list heartbeats by monitor id: %w", err)
	}
	defer rows.Close()

	heartbeats := make([]Heartbeat, 0)
	for rows.Next() {
		var heartbeat Heartbeat
		if err := rows.Scan(
			&heartbeat.ID,
			&heartbeat.MonitorID,
			&heartbeat.StatusCode,
			&heartbeat.LatencyMs,
			&heartbeat.IsUp,
			&heartbeat.ErrorMessage,
			&heartbeat.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan heartbeat: %w", err)
		}
		heartbeats = append(heartbeats, heartbeat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate heartbeats: %w", err)
	}
	return heartbeats, nil
}
