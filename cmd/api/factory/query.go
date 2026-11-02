package factory

import (
	"database/sql"

	"uptime-monitor-backend/internal/config"
	"uptime-monitor-backend/internal/query"
)

func NewQueryHandler(db *sql.DB, cfg *config.Config) *query.Handler {
	repo := query.NewRepository(db)
	generator := query.NewOpenAISQLGenerator(cfg.OpenAIAPIKey, cfg.OpenAIModel)
	embedder := query.NewOpenAIEmbedder(cfg.OpenAIAPIKey, cfg.OpenAIEmbeddingModel)
	service := query.NewService(repo, generator, embedder, cfg.QueryMaxDistance)
	return query.NewHandler(service)
}
