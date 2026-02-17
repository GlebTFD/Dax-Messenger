package main

import (
	"Messenges-service/internal/controller/websocket"
	"Messenges-service/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/hashicorp/go-hclog"
)

func main() {
	router := fiber.New()

	log := hclog.Default()

	profile := usecase.NewProfile(log)
	wc := websocket.NewWebsocketClient(profile)

	// Endpoints
	router.Get("/message", wc.MessageChanel)

	router.Listen(":8080")
}
