package usecase

import (
	"Messenges-service/internal/dto"
	"context"
	"fmt"

	"github.com/gofiber/contrib/v3/websocket"
)

func (p *Profile) MessageChannel(conn *websocket.Conn) error {
	err := conn.WriteMessage(1, []byte("Server is listening\n"))
	if err != nil {
		p.log.Error("Error to write message", "error", err)
		return fmt.Errorf("Error to write message: %s", err)
	}

	// remove context.Background to ctx from func parameters
	err = p.reader(context.Background(), conn)
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			p.log.Info("client closed the conn")
			return nil
		}

		p.log.Error("Error to read message", "error", err)
		return fmt.Errorf("read loop failed: %w", err)
	}

	return nil
}

func (p *Profile) reader(ctx context.Context, conn *websocket.Conn) error {
	for {
		var msg dto.MessageJSON

		err := conn.ReadJSON(&msg)
		if err != nil {
			return err
		}

		err = p.postgres.CreateMessage(ctx, &msg)
		if err != nil {
			p.log.Error("Error to create message", "error", err)
			// TODO: add system system_notification
		}
	}
}
