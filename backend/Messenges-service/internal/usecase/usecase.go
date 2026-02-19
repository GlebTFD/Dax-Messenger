package usecase

import (
	"github.com/hashicorp/go-hclog"
)

type Profile struct {
	log hclog.Logger
}

// maybe change name of vars
func NewProfile(l hclog.Logger) *Profile {
	return &Profile{
		log: l,
	}
}
