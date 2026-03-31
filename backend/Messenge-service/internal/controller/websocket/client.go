package websocket

import (
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/usecase"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/hashicorp/go-hclog"
)

type WebsocketClient struct {
	log            hclog.Logger
	messageService *usecase.MessageService
}

func NewWebsocketClient(p *usecase.MessageService) *WebsocketClient {
	return &WebsocketClient{
		messageService: p,
	}
}

func (wc *WebsocketClient) MessageChanel() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		err := wc.messageService.MessageChannel(c)
		if err != nil {
			wc.log.Error("Error in server", "error", err)
			return
		}
	}
}
