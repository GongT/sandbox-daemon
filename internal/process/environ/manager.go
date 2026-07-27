package environ

import (
	"os"

	internalconfig "github.com/gongt/sandbox-daemon/internal/config"
	"github.com/gongt/sandbox-daemon/internal/process/config"
)

type EnvironmentManager struct {
	config *config.EnvironmentsConfig

	initial map[string]string
}

var initialEnviron map[string]string

func NewEnvironmentManager(config *config.EnvironmentsConfig) (*EnvironmentManager, error) {
	clone := make(map[string]string)
	for k, v := range initialEnviron {
		clone[k] = v
	}

	return &EnvironmentManager{
		config:  config,
		initial: clone,
	}, nil
}

func (m *EnvironmentManager) Snapshot() internalconfig.EnvironmentMap {
	snapshot := internalconfig.EnvironmentMap{}
	snapshot.ExtendLines(os.Environ(), true)

	m.config.ApplyMap(snapshot)

	return snapshot
}
