package context

import (
	"reflect"

	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/paths"
)

type ConfigFillContext interface {
	HasValue(path paths.ConfigPath) (bool, error)
	GetArraySize(path paths.ConfigPath) (int, error)
	GetObjectKeys(path paths.ConfigPath) ([]string, error)
	GetValue(t reflect.Type, path paths.ConfigPath) (any, error)
}
