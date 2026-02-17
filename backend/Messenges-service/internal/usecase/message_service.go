package usecase

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func (p *Profile) MessageChanel() error {
	mc, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.log.Error("Error to create conn", "error", err)
		return fmt.Errorf("Error to create conn: %s", err)
	}

	err = mc.WriteMessage(1, []byte("Server is listening\n"))
	if err != nil {
		p.log.Error("Error to write message", "error", err)
		return fmt.Errorf("Error to write message: %s", err)
	}

	err = reader(mc)
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			p.log.Info("client closed the conn")
			return nil
		}

		p.log.Error("Error to read message", "error", err)
		return nil
	}

	return nil
}

func reader(conn *websocket.Conn) error {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		if string(data) == "test" {
			fmt.Printf("response: %s\n", data)
		}
	}
}
