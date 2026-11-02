package query

import (
	"context"
	"fmt"
	"strings"

	"uptime-monitor-backend/internal/auth"
)

const defaultMaxDistance = 0.15 // cosine distance; similarity >= 0.85

type Service struct {
	repo        Repository
	generator   SQLGenerator
	embedder    Embedder
	maxDistance float64
}

func NewService(repo Repository, generator SQLGenerator, embedder Embedder, maxDistance float64) *Service {
	if maxDistance <= 0 {
		maxDistance = defaultMaxDistance
	}
	return &Service{
		repo:        repo,
		generator:   generator,
		embedder:    embedder,
		maxDistance: maxDistance,
	}
}

func (s *Service) Ask(ctx context.Context, question string) (string, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("unauthorized")
	}

	question = strings.TrimSpace(question)
	if question == "" {
		return CouldNotFindAnswer, nil
	}

	embedding, err := s.embedder.Embed(ctx, question)
	if err != nil {
		return "", fmt.Errorf("embed question: %w", err)
	}

	sqlQuery, fromCache, err := s.resolveSQL(ctx, userID, question, embedding)
	if err != nil {
		return "", err
	}
	if sqlQuery == "" {
		return CouldNotFindAnswer, nil
	}

	if err := ValidateReadOnlySQL(sqlQuery); err != nil {
		if fromCache {
			// Stale/bad cache entry: regenerate once.
			sqlQuery, err = s.generateAndCache(ctx, userID, question, embedding)
			if err != nil {
				return "", err
			}
			if sqlQuery == "" {
				return CouldNotFindAnswer, nil
			}
			if err := ValidateReadOnlySQL(sqlQuery); err != nil {
				return CouldNotFindAnswer, nil
			}
		} else {
			return CouldNotFindAnswer, nil
		}
	}

	columns, rows, err := s.repo.ExecuteReadOnly(ctx, sqlQuery, userID)
	if err != nil {
		return CouldNotFindAnswer, nil
	}
	return formatAnswer(columns, rows), nil
}

func (s *Service) resolveSQL(ctx context.Context, userID, question string, embedding []float32) (sqlQuery string, fromCache bool, err error) {
	cached, found, err := s.repo.FindSimilarSQL(ctx, userID, embedding, s.maxDistance)
	if err != nil {
		return "", false, fmt.Errorf("similarity search: %w", err)
	}
	if found {
		return cached, true, nil
	}

	sqlQuery, err = s.generateAndCache(ctx, userID, question, embedding)
	if err != nil {
		return "", false, err
	}
	return sqlQuery, false, nil
}

func (s *Service) generateAndCache(ctx context.Context, userID, question string, embedding []float32) (string, error) {
	sqlQuery, err := s.generator.GenerateSQL(ctx, question)
	if err != nil {
		return "", fmt.Errorf("generate sql: %w", err)
	}
	sqlQuery = strings.TrimSpace(sqlQuery)
	if sqlQuery == "" {
		return "", nil
	}

	if err := ValidateReadOnlySQL(sqlQuery); err != nil {
		return "", nil
	}

	if err := s.repo.SaveQuestion(ctx, userID, question, sqlQuery, embedding); err != nil {
		return "", fmt.Errorf("cache question: %w", err)
	}
	return sqlQuery, nil
}
