package websocket

import (
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/usecase"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/hashicorp/go-hclog"
)

type websocketClient struct {
	log            hclog.Logger
	messageService *usecase.MessageService
}

func NewWebsocketClient(p *usecase.MessageService) *websocketClient {
	return &websocketClient{
		messageService: p,
	}
}

func (wc *websocketClient) MessageChanel() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		err := wc.messageService.MessageChannel(c)
		if err != nil {
			wc.log.Error("Error in server", "error", err)
			return
		}
	}
}
