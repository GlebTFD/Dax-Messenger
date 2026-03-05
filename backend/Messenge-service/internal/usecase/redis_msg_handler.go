package usecase

import (
	"context"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"
	"github.com/hashicorp/go-hclog"
)

type RedisMessageHandler interface {
	HandleRedisMessage(ctx context.Context, msg dto.PubSubMessage) error
}

type UserPubSubHandler struct {
	log hclog.Logger
}

func NewUserPubSubHandler(log hclog.Logger) *UserPubSubHandler {
	return &UserPubSubHandler{
		log: log,
	}
}

func (h *UserPubSubHandler) HandleRedisMessage(ctx context.Context, msg dto.PubSubMessage) error {
	h.log.Info("redis msg", "msg", msg)

	return nil
}
