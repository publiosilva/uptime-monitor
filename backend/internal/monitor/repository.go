package monitor

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(monitor Monitor) (Monitor, error)
	ListByUserID(userID string) ([]Monitor, error)
	DeleteByID(userID, monitorID string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(monitor Monitor) (Monitor, error) {
	err := r.db.QueryRow(
		`INSERT INTO monitors (user_id, name, url, method, frequency, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, name, url, method, frequency, is_active, created_at`,
		monitor.UserID,
		monitor.Name,
		monitor.URL,
		monitor.Method,
		monitor.Frequency,
		monitor.IsActive,
	).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.URL,
		&monitor.Method,
		&monitor.Frequency,
		&monitor.IsActive,
		&monitor.CreatedAt,
	)
	if err != nil {
		return Monitor{}, fmt.Errorf("create monitor: %w", err)
	}
	return monitor, nil
}

func (r *PostgresRepository) ListByUserID(userID string) ([]Monitor, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, url, method, frequency, is_active, created_at
		 FROM monitors
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list monitors: %w", err)
	}
	defer rows.Close()

	monitors := make([]Monitor, 0)
	for rows.Next() {
		var monitor Monitor
		if err := rows.Scan(
			&monitor.ID,
			&monitor.UserID,
			&monitor.Name,
			&monitor.URL,
			&monitor.Method,
			&monitor.Frequency,
			&monitor.IsActive,
			&monitor.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan monitor: %w", err)
		}
		monitors = append(monitors, monitor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate monitors: %w", err)
	}
	return monitors, nil
}

func (r *PostgresRepository) DeleteByID(userID, monitorID string) error {
	result, err := r.db.Exec(
		`DELETE FROM monitors WHERE id = $1 AND user_id = $2`,
		monitorID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete monitor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrMonitorNotFound
	}
	return nil
}
