package config

import (
	"maps"

	"github.com/gongt/sandbox-daemon/internal/tools"
)

type EnvironmentsConfig struct {
	Add       tools.EnvironmentMap `config:"environments.add"`
	Blacklist tools.Set[string]    `config:"environments.blacklist"`
	Whitelist tools.Set[string]    `config:"environments.whitelist"`
}

func (m *EnvironmentsConfig) Set(name, value string) {
	m.Add[name] = value
	delete(m.Blacklist, name)
}

func (m *EnvironmentsConfig) Unset(name string) {
	delete(m.Add, name)
}

func (m *EnvironmentsConfig) AddBlacklist(name string) {
	m.Blacklist[name] = true
}

func (m *EnvironmentsConfig) DeleteBlacklist(name string) {
	delete(m.Blacklist, name)
}

func (m *EnvironmentsConfig) ClearBlacklist() {
	m.Blacklist = tools.Set[string]{}
}

func (m *EnvironmentsConfig) AddWhitelist(name string) {
	m.Whitelist[name] = true
}

func (m *EnvironmentsConfig) DeleteWhitelist(name string) {
	delete(m.Whitelist, name)
}

func (m *EnvironmentsConfig) ClearWhitelist() {
	m.Whitelist = tools.Set[string]{}
}

func (m *EnvironmentsConfig) ApplyMap(snapshot map[string]string) {
	for k := range m.Blacklist {
		delete(snapshot, k)
	}
	if len(m.Whitelist) > 0 {
		for k := range snapshot {
			if !m.Whitelist.Has(k) {
				delete(snapshot, k)
			}
		}
	}
	maps.Copy(snapshot, m.Add)
}
