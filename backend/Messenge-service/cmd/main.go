package main

import (
	"context"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/config"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/postgres"
	wsClient "github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/controller/websocket"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/usecase"

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
