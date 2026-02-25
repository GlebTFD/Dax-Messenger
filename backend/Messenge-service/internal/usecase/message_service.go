package usecase

import (
	"context"
	"fmt"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"

	"github.com/gofiber/contrib/v3/websocket"
)

func (m *MessageService) MessageChannel(conn *websocket.Conn) error {
	err := conn.WriteMessage(1, []byte("Server is listening\n"))
	if err != nil {
		m.log.Error("Error to write message", "error", err)
		return fmt.Errorf("Error to write message: %s", err)
	}

	// remove context.Background to ctx from func parameters
	err = m.wsReader(context.Background(), conn)
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			m.log.Info("client closed the conn")
			return nil
		}

		m.log.Error("Error to read message", "error", err)
		return fmt.Errorf("read loop failed: %w", err)
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
	}
}

func (m *MessageService) subscribeToNewMessage() error {
	return nil
}
