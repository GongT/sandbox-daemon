package config_loader

import (
	"github.com/gongt/sandbox-daemon/packages/config_loader/internal"
	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/context"
)

type ConfigFillContext = context.ConfigFillContext

func WalkStruct(input any, ctx context.ConfigFillContext) error {
	return internal.WalkStruct(input, ctx)
}
