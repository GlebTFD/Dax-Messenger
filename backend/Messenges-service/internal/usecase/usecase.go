package usecase

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-hclog"
)

type Profile struct {
	log      hclog.Logger
	upgrader websocket.Upgrader
}

// maybe change name of vars
func NewProfile(l hclog.Logger) *Profile {
	var u = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return &Profile{
		log:      l,
		upgrader: u,
	}
}
