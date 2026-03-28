package usecase

import (
	"context"
	"fmt"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/domain"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"

	"github.com/gofiber/contrib/v3/websocket"
)

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

	errChan := make(chan error, 2)

	go func() {
		// CONTEXT.BACKGROUD!!!
		err := m.wsReader(context.Background(), conn)
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			errChan <- fmt.Errorf("read loop failed: %w", err)
		} else {
			errChan <- nil
		}
	}()

	// maybe this needs to be moved to the controller
	go func() {
		// CONTEXT.BACKGROUD!!!
		err := m.redisPubSub.SubscribeAndRun(context.Background(), "chat:"+id.ID)
		if err != nil {
			errChan <- fmt.Errorf("subscribe failed: %w", err)
		} else {
			errChan <- nil
		}
	}()

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func (m *MessageService) wsReader(ctx context.Context, conn *websocket.Conn) error {
	for {
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

		err = m.redisPubSub.PublishToChannel(ctx, "chat:"+msg.Payload.ReplyTo, msg.Payload.Text)
		if err != nil {
			m.log.Error("Error to publish msg to channel", "error", err)
		}
	}
}
