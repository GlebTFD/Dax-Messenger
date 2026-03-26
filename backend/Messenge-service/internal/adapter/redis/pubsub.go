package redis

import (
	"context"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/usecase"
	"github.com/hashicorp/go-hclog"
	"github.com/redis/go-redis/v9"
)

type PubSubConfig struct {
	Address  string
	Password string
	DB       int
	Channel  string
}

type RedisPubSubClient struct {
	log     hclog.Logger
	client  *redis.Client
	channel string
	handler usecase.RedisMessageHandler
}

func New(ctx context.Context, log hclog.Logger, cfg PubSubConfig, handler usecase.RedisMessageHandler) *RedisPubSubClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:        cfg.Address,
		Password:    cfg.Password,
		DB:          cfg.DB,
		ReadTimeout: -1,
		MaxRetries:  -1,
	})

	return &RedisPubSubClient{
		log:     log,
		client:  rdb,
		channel: cfg.Channel,
		handler: handler,
	}
}

func (p *RedisPubSubClient) SubscribeAndRun(ctx context.Context, channelName string) error {
	pubsub := p.client.Subscribe(ctx, channelName)
	defer pubsub.Close()

	ch := pubsub.Channel()
	p.log.Info("Subscribed to channel", "channel", channelName)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			if p.handler != nil {
				if err := p.handler.HandleRedisMessage(ctx, dto.PubSubMessage{
					Channel: msg.Channel,
					Payload: msg.Payload,
				}); err != nil {
					p.log.Error("Handler Error", "error", err, "channel", msg.Channel)
				}
			}
		}
	}
}

func (p *RedisPubSubClient) PublishToChannel(ctx context.Context, channel string, msg interface{}) error {
	return p.client.Publish(ctx, channel, msg).Err()
}
