package usecase

import (
	"context"
	"fmt"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"
)

type RedisMessageHandler interface {
	HandleRedisMessage(ctx context.Context, msg dto.RedisMessage) error
}

func (h *UserPubSubHandler) HandleRedisMessage(ctx context.Context, msg dto.RedisMessage) error {
	h.log.Info("redis msg", "type", msg.Type, "channel", msg.Channel)

	id := msg.Channel[len("chat:"):]

	h.wsConns.RLock()
	conn, ok := h.wsConns.Conns[id]
	h.wsConns.RUnlock()

	if !ok {
		h.log.Warn("User not connected, message dropped", "user_id", id)
		return nil
	}

	var err error
	switch msg.Type {
	case "message", "message_deleted", "message_updated":
		err = conn.WriteJSON(msg)
	default:
		h.log.Warn("Unknown message type from redis", "type", msg.Type)
		return nil
	}

	if err != nil {
		h.log.Error("Error to send msg to ws", "error", err)
		return fmt.Errorf("error to send msg to ws: %w", err)
	}

	return nil
}
