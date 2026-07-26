package initializer

import "github.com/gongt/sandbox-daemon/internal/config"

type Initializer struct {
	config *config.Config
}

func New(cfg *config.Config) *Initializer {
	return &Initializer{config: cfg}
}

func (i *Initializer) Run() {
}

func (i *Initializer) Destroy() {
}
