package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/domain"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/jackc/pgx/v5"
)

var ErrMessageNotFound = errors.New("message not found")
var ErrEmptyText = errors.New("text cannot be empty")

// http
func (m *MessageService) DeleteMessage(msgId string) error {
	replyTo, err := m.postgres.DeleteMessage(context.Background(), msgId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}

	notification := dto.RedisMessage{
		Type:    "message_deleted",
		Payload: dto.DeletedNotificationPayload{ID: msgId},
	}
	if err := m.redisPubSub.PublishToChannel(context.Background(), "chat:"+replyTo, notification); err != nil {
		m.log.Error("Error to publish delete notification", "error", err)
	}

	return nil
}

func (m *MessageService) UpdateMessage(msgId string, text string) error {
	if text == "" {
		return ErrEmptyText
	}

	replyTo, err := m.postgres.UpdateMessage(context.Background(), msgId, text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}

	notification := dto.RedisMessage{
		Type:    "message_updated",
		Payload: dto.UpdatedNotificationPayload{ID: msgId, Text: text},
	}
	if err := m.redisPubSub.PublishToChannel(context.Background(), "chat:"+replyTo, notification); err != nil {
		m.log.Error("Error to publish update notification", "error", err)
	}

	return nil
}

// websocket
func (m *MessageService) MessageChannel(conn *websocket.Conn) error {
	// user wiil be send id
	var id domain.UserId
	err := conn.ReadJSON(&id)
	if err != nil {
		m.log.Error("User didt send id ot read json fatal error", "error", err)
		return fmt.Errorf("user didt send id ot read json fatal error")
	}

	// add to ws wsConns
	m.wsConns.Lock()
	m.wsConns.Conns[id.ID] = conn
	m.wsConns.Unlock()

	// cleanup on disconnect
	defer func() {
		m.wsConns.Lock()
		delete(m.wsConns.Conns, id.ID)
		m.wsConns.Unlock()
	}()

	errChan := make(chan error, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := m.wsReader(ctx, conn)
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			errChan <- fmt.Errorf("read loop failed: %w", err)
		} else {
			errChan <- nil
		}
	}()

	go func() {
		err := m.redisPubSub.SubscribeAndRun(ctx, "chat:"+id.ID)
		if err != nil {
			errChan <- fmt.Errorf("subscribe failed: %w", err)
		} else {
			errChan <- nil
		}
	}()

	// Wait for both goroutines to finish. Cancel context on first error to stop the other.
	var firstErr error
	for range 2 {
		if err := <-errChan; err != nil && firstErr == nil {
			firstErr = err
			cancel()
		}
	}

	return firstErr
}

func (m *MessageService) wsReader(ctx context.Context, conn *websocket.Conn) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var msg dto.MessageJSON

		err := conn.ReadJSON(&msg)
		if err != nil {
			return err
		}

		err = m.postgres.CreateMessage(ctx, &msg)
		if err != nil {
			m.log.Error("Error to create message", "error", err)
			// TODO: add system system_notification
		}

		redisMsg := dto.RedisMessage{
			Type:    "message",
			Payload: msg.Payload,
		}
		if err := m.redisPubSub.PublishToChannel(ctx, "chat:"+msg.Payload.ReplyTo, redisMsg); err != nil {
			m.log.Error("Error to publish msg to channel", "error", err)
		}
	}
}
