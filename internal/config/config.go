package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL           string
	APIPort               string
	JWTSecret             string
	WSPort                string
	OpenAIAPIKey          string
	OpenAIModel           string
	OpenAIEmbeddingModel  string
	QueryMaxDistance      float64
}

func Load() (Config, error) {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok || dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok || jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3333"
	}

	wsPort := os.Getenv("WS_PORT")
	if wsPort == "" {
		wsPort = "3334"
	}

	openAIModel := os.Getenv("OPENAI_MODEL")
	if openAIModel == "" {
		openAIModel = "gpt-4o-mini"
	}

	embeddingModel := os.Getenv("OPENAI_EMBEDDING_MODEL")
	if embeddingModel == "" {
		embeddingModel = "text-embedding-3-small"
	}

	maxDistance := 0.15
	if raw := os.Getenv("QUERY_MAX_DISTANCE"); raw != "" {
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return Config{}, fmt.Errorf("QUERY_MAX_DISTANCE: %w", err)
		}
		maxDistance = parsed
	}

	return Config{
		DatabaseURL:          dbURL,
		APIPort:              port,
		JWTSecret:            jwtSecret,
		WSPort:               wsPort,
		OpenAIAPIKey:         os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:          openAIModel,
		OpenAIEmbeddingModel: embeddingModel,
		QueryMaxDistance:     maxDistance,
	}, nil
}
