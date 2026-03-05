package config

import (
	"fmt"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/postgres"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/redis"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres postgres.Config
	PubSub   redis.PubSubConfig
}

func InitConfig() (Config, error) {
	cfg := Config{}

	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error to load env: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("error to parse env: %w", err)
	}

	return cfg, nil
}
