package usecase

import (
	"context"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/postgres"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/domain"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"

	"github.com/hashicorp/go-hclog"
)

type Postgres interface {
	CreateMessage(ctx context.Context, msg *dto.MessageJSON) error
	DeleteMessage(ctx context.Context, msgId string) (string, error)
	UpdateMessage(ctx context.Context, msgId string, text string) (string, error)
}

type RedisPubSub interface {
	SubscribeAndRun(ctx context.Context, channelName string) error
	PublishToChannel(ctx context.Context, channel string, msg dto.RedisMessage) error
}

// MessageService
type MessageService struct {
	log         hclog.Logger
	postgres    Postgres
	redisPubSub RedisPubSub
	wsConns     *domain.ConnectionManager
}

func NewMessageService(log hclog.Logger, postgres *postgres.Pool, redisPubSub RedisPubSub, wsConns *domain.ConnectionManager) *MessageService {
	return &MessageService{
		log:         log,
		postgres:    postgres,
		redisPubSub: redisPubSub,
		wsConns:     wsConns,
	}
}

// UserPubSubHandler
type UserPubSubHandler struct {
	log     hclog.Logger
	wsConns *domain.ConnectionManager
}

func NewUserPubSubHandler(log hclog.Logger, wsConns *domain.ConnectionManager) *UserPubSubHandler {
	return &UserPubSubHandler{
		log:     log,
		wsConns: wsConns,
	}
}
