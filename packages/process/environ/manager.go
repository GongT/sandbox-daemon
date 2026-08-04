package environ

import (
	"os"
)

type EnvironmentManager struct {
	config *ManagerConfig

	initial map[string]string
}

func NewManager(config *ManagerConfig) *EnvironmentManager {
	snapshot := Map{}
	snapshot.ExtendLines(os.Environ(), true)

	return &EnvironmentManager{
		config:  config,
		initial: snapshot,
	}
}

func (m *EnvironmentManager) Snapshot() (Map, error) {
	snapshot := Map{}
	snapshot.Extend(m.initial, true)

	err := m.config.ApplyMap(snapshot)

	return snapshot, err
}
