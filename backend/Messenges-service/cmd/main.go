package main

import (
	"Messenges-service/config"
	"Messenges-service/internal/adapter/postgres"
	wsClient "Messenges-service/internal/controller/websocket"
	"Messenges-service/internal/usecase"
	"context"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/hashicorp/go-hclog"
)

func main() {
	ctx := context.Background()

	router := fiber.New()

	log := hclog.Default()

	config, err := config.InitConfig()
	if err != nil {
		log.Error("Error to init config", "error", err)
		return
	}

	pool, err := postgres.New(ctx, log, config.Postgres)
	if err != nil {
		log.Error("Error to init pool to db", "error", err)
		return
	}

	profile := usecase.NewProfile(log, pool)
	wc := wsClient.NewWebsocketClient(profile)

	// Endpoints
	router.Get("/message", websocket.New(wc.MessageChanel()))

	router.Listen(":8080")
}
