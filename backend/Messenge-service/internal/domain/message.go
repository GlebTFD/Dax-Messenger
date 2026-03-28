package domain

import (
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

type UserId struct {
	ID string `json:"id"`
}

type ConnectionManager struct {
	sync.RWMutex
	Conns map[string]*websocket.Conn
}
