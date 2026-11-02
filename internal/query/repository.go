package query

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pgvector/pgvector-go"
)

// Repository persists question embeddings and executes guarded read-only SQL.
type Repository interface {
	FindSimilarSQL(ctx context.Context, userID string, embedding []float32, maxDistance float64) (sqlQuery string, found bool, err error)
	SaveQuestion(ctx context.Context, userID, question, sqlQuery string, embedding []float32) error
	ExecuteReadOnly(ctx context.Context, sqlQuery string, args ...any) (columns []string, rows [][]string, err error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindSimilarSQL(ctx context.Context, userID string, embedding []float32, maxDistance float64) (string, bool, error) {
	var sqlQuery string
	err := r.db.QueryRowContext(ctx, `
		SELECT sql_query
		FROM query_question_cache
		WHERE user_id = $1
		  AND (embedding <=> $2) <= $3
		ORDER BY embedding <=> $2
		LIMIT 1
	`, userID, pgvector.NewVector(embedding), maxDistance).Scan(&sqlQuery)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("find similar sql: %w", err)
	}
	return sqlQuery, true, nil
}

func (r *PostgresRepository) SaveQuestion(ctx context.Context, userID, question, sqlQuery string, embedding []float32) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO query_question_cache (user_id, question, sql_query, embedding)
		VALUES ($1, $2, $3, $4)
	`, userID, question, sqlQuery, pgvector.NewVector(embedding))
	if err != nil {
		return fmt.Errorf("save question cache: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ExecuteReadOnly(ctx context.Context, sqlQuery string, args ...any) ([]string, [][]string, error) {
	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("execute read-only sql: %w", err)
	}
	defer rows.Close()

	colTypes, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("columns: %w", err)
	}

	results := make([][]string, 0)
	for rows.Next() {
		raw := make([]any, len(colTypes))
		ptrs := make([]any, len(colTypes))
		for i := range raw {
			ptrs[i] = &raw[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, fmt.Errorf("scan row: %w", err)
		}
		row := make([]string, len(colTypes))
		for i, v := range raw {
			row[i] = stringifyDBValue(v)
		}
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate rows: %w", err)
	}
	return colTypes, results, nil
}

func stringifyDBValue(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case []byte:
		return string(t)
	default:
		return fmt.Sprint(t)
	}
}

func formatAnswer(columns []string, rows [][]string) string {
	if len(rows) == 0 {
		return CouldNotFindAnswer
	}

	var b strings.Builder
	for i, row := range rows {
		if i > 0 {
			b.WriteString("; ")
		}
		for j, col := range columns {
			if j > 0 {
				b.WriteString(", ")
			}
			b.WriteString(col)
			b.WriteString(": ")
			if j < len(row) {
				b.WriteString(row[j])
			}
		}
	}
	return b.String()
}
