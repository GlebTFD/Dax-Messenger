package main

import (
	wsClient "Messenges-service/internal/controller/websocket"
	"Messenges-service/internal/usecase"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/hashicorp/go-hclog"
)

func main() {
	router := fiber.New()

	log := hclog.Default()

	profile := usecase.NewProfile(log)
	wc := wsClient.NewWebsocketClient(profile)

	// Endpoints
	router.Get("/message", websocket.New(wc.MessageChanel()))

	router.Listen(":8080")
}
