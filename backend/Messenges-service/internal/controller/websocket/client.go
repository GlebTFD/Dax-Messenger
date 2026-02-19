package websocket

import (
	"Messenges-service/internal/usecase"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/hashicorp/go-hclog"
)

type websocketClient struct {
	log     hclog.Logger
	profile *usecase.Profile
}

func NewWebsocketClient(p *usecase.Profile) *websocketClient {
	return &websocketClient{
		profile: p,
	}
}

func (wc *websocketClient) MessageChanel() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		_ = wc.profile.MessageChanel(c)
	}
}
