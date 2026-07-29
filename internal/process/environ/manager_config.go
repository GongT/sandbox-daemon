package environ

import (
	"maps"

	"github.com/gongt/sandbox-daemon/internal/tools/types"
	"github.com/pkg/errors"
)

type ManagerConfig struct {
	Add       Map               `config:"environments.add"`
	Blacklist types.Set[string] `config:"environments.blacklist"`
	Whitelist types.Set[string] `config:"environments.whitelist"`
	Requires  types.Set[string] `config:"environments.requires"`
}

func (m *ManagerConfig) Set(name, value string) {
	m.Add[name] = value
}

func (m *ManagerConfig) Unset(name string) {
	m.Add.Delete(name)
}

func (m *ManagerConfig) AddBlacklist(name string) {
	m.Blacklist.Add(name)
}

func (m *ManagerConfig) DeleteBlacklist(name string) {
	m.Blacklist.Delete(name)
}

func (m *ManagerConfig) ClearBlacklist() {
	m.Blacklist = types.Set[string]{}
}

func (m *ManagerConfig) AddWhitelist(name string) {
	m.Whitelist.Add(name)
}

func (m *ManagerConfig) DeleteWhitelist(name string) {
	m.Whitelist.Delete(name)
}

func (m *ManagerConfig) ClearWhitelist() {
	m.Whitelist = types.Set[string]{}
}

func (m *ManagerConfig) ApplyMap(snapshot map[string]string) error {
	for _, k := range m.Blacklist {
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
	for _, k := range m.Requires {
		if _, ok := snapshot[k]; !ok {
			return errors.Errorf("必须存在的环境变量 %q 不存在", k)
		}
	}
	return nil
}
