package redis

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/redis/go-redis/v9"
)

type PubSubConfig struct {
	Address  string
	Password string
	DB       int
}

type RedisPubSubClient struct {
	log    hclog.Logger
	client *redis.Client
}

func (rc *RedisPubSubClient) New(ctx context.Context, log hclog.Logger, cfg PubSubConfig) *RedisPubSubClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:        cfg.Address,
		Password:    cfg.Password,
		DB:          cfg.DB,
		ReadTimeout: -1,
		MaxRetries:  -1,
	})

	return &RedisPubSubClient{
		log:    log,
		client: rdb,
	}
}
