package environ

import (
	"maps"
	"os"

	"github.com/gongt/sandbox-daemon/internal/config"
)

type EnvironmentManager struct {
	add       config.EnvironmentMap `config:"environments.add"`
	blacklist config.Set[string]    `config:"environments.blacklist"`
	whitelist config.Set[string]    `config:"environments.whitelist"`

	initial map[string]string
}

var initialEnviron map[string]string

func NewEnvironmentManager() (*EnvironmentManager, error) {
	clone := make(map[string]string)
	for k, v := range initialEnviron {
		clone[k] = v
	}

	return &EnvironmentManager{
		add:       make(config.EnvironmentMap),
		blacklist: config.Set[string]{"PWD": true},
		whitelist: config.Set[string]{},
		initial:   clone,
	}, nil
}

func (m *EnvironmentManager) Set(name, value string) {
	m.add[name] = value
	delete(m.blacklist, name)
}

func (m *EnvironmentManager) Unset(name string) {
	delete(m.add, name)
}

func (m *EnvironmentManager) Blacklist(name string) {
	m.blacklist[name] = true
}

func (m *EnvironmentManager) UnBlacklist(name string) {
	delete(m.blacklist, name)
}

func (m *EnvironmentManager) ClearBlacklist() {
	m.blacklist = config.Set[string]{}
}

func (m *EnvironmentManager) Whitelist(name string) {
	m.whitelist[name] = true
}

func (m *EnvironmentManager) UnWhitelist(name string) {
	delete(m.whitelist, name)
}

func (m *EnvironmentManager) ClearWhitelist() {
	m.whitelist = config.Set[string]{}
}

func (m *EnvironmentManager) Snapshot() config.EnvironmentMap {
	snapshot := config.EnvironmentMap{}
	snapshot.ExtendLines(os.Environ(), true)

	maps.Copy(snapshot, m.add)

	for k := range m.blacklist {
		delete(snapshot, k)
	}
	if len(m.whitelist) > 0 {
		for k := range snapshot {
			if !m.whitelist.Has(k) {
				delete(snapshot, k)
			}
		}
	}
	return snapshot
}
