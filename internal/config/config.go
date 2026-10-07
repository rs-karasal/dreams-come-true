package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	return Config{
		Addr: envOrDefault("HTTP_ADDR", ":8080"),
	}, nil
}

func envOrDefault(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
