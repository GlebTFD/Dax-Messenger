package websocket

import (
	"Messenges-service/internal/usecase"

	"github.com/gofiber/fiber/v3"
)

type websocketClient struct {
	profile *usecase.Profile
}

func NewWebsocketClient(p *usecase.Profile) *websocketClient {
	return &websocketClient{
		profile: p,
	}
}

func (wc *websocketClient) MessageChanel(c fiber.Ctx) {
	err := wc.profile.MessageChanel()
}
