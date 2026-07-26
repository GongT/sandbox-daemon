package daemon

import "github.com/gongt/sandbox-daemon/internal/config"

type Daemon struct {
	config *config.Config
}

func New(cfg *config.Config) *Daemon {
	return &Daemon{config: cfg}
}

func (d *Daemon) Run() {
}

func (d *Daemon) Destroy() {
}
