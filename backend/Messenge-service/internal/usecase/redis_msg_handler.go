package usecase

import (
	"context"
	"fmt"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"
)

type RedisMessageHandler interface {
	HandleRedisMessage(ctx context.Context, msg dto.PubSubMessage) error
}

func (h *UserPubSubHandler) HandleRedisMessage(ctx context.Context, msg dto.PubSubMessage) error {
	h.log.Info("redis msg", "msg", msg)

	// Make it so that you can get the ID in a different way, if necessary
	id := msg.Channel[len("chat:"):]

	payload := dto.TextMessagePayload{
		ReplyTo: id,
		Text:    msg.Payload,
	}
	msgJSON := &dto.MessageJSON{
		Type:    "message",
		Payload: payload,
	}

	h.wsConns.RLock()
	conn, ok := h.wsConns.Conns[id]
	h.wsConns.RUnlock()

	if !ok {
		h.log.Warn("User not connected, message dropped", "user_id", id)
		return nil
	}

	err := conn.WriteJSON(msgJSON)
	if err != nil {
		h.log.Error("Error to send new msg", "error", err)
		return fmt.Errorf("error to send new msg")
	}

	return nil
}
