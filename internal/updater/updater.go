package updater

import "github.com/gongt/sandbox-daemon/internal/config"

type Updater struct {
	config *config.Config
}

func New(cfg *config.Config) *Updater {
	return &Updater{config: cfg}
}

func (u *Updater) Run() {
}

func (u *Updater) Destroy() {
}
