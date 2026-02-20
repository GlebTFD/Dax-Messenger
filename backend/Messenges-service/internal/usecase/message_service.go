package usecase

import (
	"Messenges-service/internal/dto"
	"fmt"

	"github.com/gofiber/contrib/v3/websocket"
)

func (p *Profile) MessageChanel(conn *websocket.Conn) error {
	err := conn.WriteMessage(1, []byte("Server is listening\n"))
	if err != nil {
		p.log.Error("Error to write message", "error", err)
		return fmt.Errorf("Error to write message: %s", err)
	}

	err = reader(conn)
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

func reader(conn *websocket.Conn) error {
	var msg dto.MessageJSON

	for {
		err := conn.ReadJSON(&msg)
		if err != nil {
			return err
		}

		// remove if not needed
		conn.WriteMessage(1, []byte("msg is very cool"))

		// for test
		fmt.Printf("msg: %s\n", msg.Payload.Text)
	}
}
