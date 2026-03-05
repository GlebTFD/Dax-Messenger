package usecase

import (
	"context"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/postgres"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"

	"github.com/hashicorp/go-hclog"
)

type Postgres interface {
	CreateMessage(ctx context.Context, msg *dto.MessageJSON) error
}

type RedisPubSub interface {
	SubscribeAndRun(ctx context.Context, channelName string) error
	PublishToChannel(ctx context.Context, channel string, msg []byte) error
}

type MessageService struct {
	log         hclog.Logger
	postgres    Postgres
	redisPubSub RedisPubSub
}

// maybe change name of vars
func NewMessageService(log hclog.Logger, postgres *postgres.Pool, redisPubSub RedisPubSub) *MessageService {
	return &MessageService{
		log:         log,
		postgres:    postgres,
		redisPubSub: redisPubSub,
	}
}
