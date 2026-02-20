package usecase

import (
	"Messenges-service/internal/adapter/postgres"
	"Messenges-service/internal/dto"
	"context"

	"github.com/hashicorp/go-hclog"
)

type Postgres interface {
	CreateMessage(ctx context.Context, msg dto.Message) error
}

type Profile struct {
	log      hclog.Logger
	postgres Postgres
}

// maybe change name of vars
func NewProfile(l hclog.Logger, postgres *postgres.Pool) *Profile {
	return &Profile{
		log:      l,
		postgres: postgres,
	}
}
