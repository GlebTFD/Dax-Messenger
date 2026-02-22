package usecase

import (
	"context"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/adapter/postgres"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"

	"github.com/hashicorp/go-hclog"
)

type Postgres interface {
	CreateMessage(ctx context.Context, msg *dto.MessageJSON) error
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
