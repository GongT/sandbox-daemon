package environ

import (
	"os"

	"github.com/gongt/sandbox-daemon/internal/process/config"
	"github.com/gongt/sandbox-daemon/internal/tools"
)

type EnvironmentManager struct {
	config *config.EnvironmentsConfig

	initial map[string]string
}

var initialEnviron map[string]string

func New(config *config.EnvironmentsConfig) *EnvironmentManager {
	clone := make(map[string]string)
	for k, v := range initialEnviron {
		clone[k] = v
	}

	return &EnvironmentManager{
		config:  config,
		initial: clone,
	}
}

func (m *EnvironmentManager) Snapshot() tools.EnvironmentMap {
	snapshot := tools.EnvironmentMap{}
	snapshot.ExtendLines(os.Environ(), true)

	m.config.ApplyMap(snapshot)

	return snapshot
}
