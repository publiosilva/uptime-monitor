package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	APIPort     string
	JWTSecret   string
	WSPort      string
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

	return Config{
		DatabaseURL: dbURL,
		APIPort:     port,
		JWTSecret:   jwtSecret,
		WSPort:      wsPort,
	}, nil
}
