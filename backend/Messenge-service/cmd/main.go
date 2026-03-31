package main

import (
	"context"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/config"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/postgres"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/redis"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/controller/http"
	wsClient "github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/controller/websocket"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/domain"
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

	// add delete conn after disconn
	wsConns := domain.ConnectionManager{
		Conns: make(map[string]*websocket.Conn),
	}

	redisPSHandler := usecase.NewUserPubSubHandler(log, &wsConns)
	psRedisClient := redis.New(ctx, log, config.PubSub, redisPSHandler)

	messageService := usecase.NewMessageService(log, pool, psRedisClient, &wsConns)
	wc := wsClient.NewWebsocketClient(messageService)
	http := http.NewHTTPHandler(log, messageService)

	// Endpoints
	router.Get("/message", websocket.New(wc.MessageChanel()))
	router.Delete("/message/:id", http.DeleteMessage)

	err = router.Listen(":8080")
	if err != nil {
		log.Error("error in listen server", "error", err)
	}
}
