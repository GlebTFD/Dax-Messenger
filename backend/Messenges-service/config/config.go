package config

import (
	"Messenges-service/internal/adapter/postgres"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres postgres.Config
}

func InitConfig() (Config, error) {
	cfg := Config{}

	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("Error to load env: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("Error to parse env: %w", err)
	}

	return cfg, nil
}
