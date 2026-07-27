package config

import (
	"maps"

	"github.com/gongt/sandbox-daemon/internal/tools"
)

type EnvironmentsConfig struct {
	add       tools.EnvironmentMap `config:"environments.add"`
	blacklist tools.Set[string]    `config:"environments.blacklist"`
	whitelist tools.Set[string]    `config:"environments.whitelist"`
}

func (m *EnvironmentsConfig) Set(name, value string) {
	m.add[name] = value
	delete(m.blacklist, name)
}

func (m *EnvironmentsConfig) Unset(name string) {
	delete(m.add, name)
}

func (m *EnvironmentsConfig) Blacklist(name string) {
	m.blacklist[name] = true
}

func (m *EnvironmentsConfig) UnBlacklist(name string) {
	delete(m.blacklist, name)
}

func (m *EnvironmentsConfig) ClearBlacklist() {
	m.blacklist = tools.Set[string]{}
}

func (m *EnvironmentsConfig) Whitelist(name string) {
	m.whitelist[name] = true
}

func (m *EnvironmentsConfig) UnWhitelist(name string) {
	delete(m.whitelist, name)
}

func (m *EnvironmentsConfig) ClearWhitelist() {
	m.whitelist = tools.Set[string]{}
}

func (m *EnvironmentsConfig) ApplyMap(snapshot map[string]string) {
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
}
