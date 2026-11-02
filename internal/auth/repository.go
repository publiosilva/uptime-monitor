package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	CreateUser(email, passwordHash string) (User, error)
	GetUserByEmail(email string) (User, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(email, passwordHash string) (User, error) {
	var user User
	err := r.db.QueryRow(
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at`,
		email,
		passwordHash,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailTaken
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (r *PostgresRepository) GetUserByEmail(email string) (User, error) {
	var user User
	err := r.db.QueryRow(
		`SELECT id, email, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	switch {
	case err == sql.ErrNoRows:
		return User{}, ErrUserNotFound
	case err != nil:
		return User{}, fmt.Errorf("query row: %w", err)
	default:
		return user, nil
	}
}
