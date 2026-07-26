package logforwarder

import "github.com/gongt/sandbox-daemon/internal/config"

type LogForwarder struct {
	config *config.Config
}

func New(cfg *config.Config) *LogForwarder {
	return &LogForwarder{config: cfg}
}

func (l *LogForwarder) Run() {
}

func (l *LogForwarder) Destroy() {
}
