package monitor

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository interface {
	Create(monitor Monitor) (Monitor, error)
	ListByUserID(userID string) ([]Monitor, error)
	ListAll() ([]Monitor, error)
	DeleteByID(userID, monitor_id string) error
	FindByUserAndID(userID, monitor_id string) (Monitor, error)
	FindByID(monitor_id string) (Monitor, error)
	Update(monitor Monitor) error
	UpdateByUserAndID(userID string, monitor Monitor) (Monitor, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(monitor Monitor) (Monitor, error) {
	err := r.db.QueryRow(
		`INSERT INTO monitors (user_id, name, url, method, timeout, frequency, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, user_id, name, url, method, timeout, frequency, is_active, is_up, created_at`,
		monitor.UserID,
		monitor.Name,
		monitor.URL,
		monitor.Method,
		monitor.Timeout,
		monitor.Frequency,
		monitor.IsActive,
	).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.URL,
		&monitor.Method,
		&monitor.Timeout,
		&monitor.Frequency,
		&monitor.IsActive,
		&monitor.IsUp,
		&monitor.CreatedAt,
	)
	if err != nil {
		return Monitor{}, fmt.Errorf("create monitor: %w", err)
	}
	return monitor, nil
}

func (r *PostgresRepository) ListByUserID(userID string) ([]Monitor, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, url, method, timeout, frequency, is_active, is_up, created_at
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
			&monitor.Timeout,
			&monitor.Frequency,
			&monitor.IsActive,
			&monitor.IsUp,
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

func (r *PostgresRepository) ListAll() ([]Monitor, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, url, method, timeout, frequency, is_active, is_up, created_at
		 FROM monitors
		 ORDER BY created_at DESC`,
	)

	if err != nil {
		return nil, fmt.Errorf("list all monitors: %w", err)
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
			&monitor.Timeout,
			&monitor.Frequency,
			&monitor.IsActive,
			&monitor.IsUp,
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

func (r *PostgresRepository) DeleteByID(userID, monitor_id string) error {
	result, err := r.db.Exec(
		`DELETE FROM monitors WHERE id = $1 AND user_id = $2`,
		monitor_id,
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

func (r *PostgresRepository) FindByUserAndID(userID, monitor_id string) (Monitor, error) {
	var monitor Monitor
	err := r.db.QueryRow(
		`SELECT id, user_id, name, url, method, timeout, frequency, is_active, is_up, created_at
		FROM monitors
		WHERE id = $1 AND user_id = $2`,
		monitor_id,
		userID,
	).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.URL,
		&monitor.Method,
		&monitor.Timeout,
		&monitor.Frequency,
		&monitor.IsActive,
		&monitor.IsUp,
		&monitor.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Monitor{}, ErrMonitorNotFound
		}
		return Monitor{}, fmt.Errorf("find monitor by user and id: %w", err)
	}
	return monitor, nil
}

func (r *PostgresRepository) FindByID(monitor_id string) (Monitor, error) {
	var monitor Monitor
	err := r.db.QueryRow(
		`SELECT id, user_id, name, url, method, timeout, frequency, is_active, is_up, created_at
		FROM monitors
		WHERE id = $1`,
		monitor_id,
	).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.URL,
		&monitor.Method,
		&monitor.Timeout,
		&monitor.Frequency,
		&monitor.IsActive,
		&monitor.IsUp,
		&monitor.CreatedAt,
	)
	if err != nil {
		return Monitor{}, fmt.Errorf("find monitor by id: %w", err)
	}
	return monitor, nil
}

func (r *PostgresRepository) Update(monitor Monitor) error {
	result, err := r.db.Exec(
		`UPDATE monitors SET is_up = $1 WHERE id = $2`,
		monitor.IsUp,
		monitor.ID,
	)
	if err != nil {
		return fmt.Errorf("update monitor: %w", err)
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

func (r *PostgresRepository) UpdateByUserAndID(userID string, monitor Monitor) (Monitor, error) {
	err := r.db.QueryRow(
		`UPDATE monitors
		 SET name = $1, url = $2, method = $3, timeout = $4, frequency = $5, is_active = $6
		 WHERE id = $7 AND user_id = $8
		 RETURNING id, user_id, name, url, method, timeout, frequency, is_active, is_up, created_at`,
		monitor.Name,
		monitor.URL,
		monitor.Method,
		monitor.Timeout,
		monitor.Frequency,
		monitor.IsActive,
		monitor.ID,
		userID,
	).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.URL,
		&monitor.Method,
		&monitor.Timeout,
		&monitor.Frequency,
		&monitor.IsActive,
		&monitor.IsUp,
		&monitor.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Monitor{}, ErrMonitorNotFound
		}
		return Monitor{}, fmt.Errorf("update monitor by user and id: %w", err)
	}
	return monitor, nil
}
