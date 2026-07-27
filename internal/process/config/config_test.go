package config

import (
	_ "embed"
	"log"
	"testing"

	"github.com/goforj/godump"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	assets "github.com/gongt/sandbox-daemon"
	internalconfig "github.com/gongt/sandbox-daemon/internal/config"
)

func TestProcessConfigLoad(t *testing.T) {
	log.SetOutput(t.Output())

	cfg := LifecycleConfig{}
	err := internalconfig.LoadConfigContent(assets.ExampleConfigYAML, &cfg.exec, &cfg.stop, &cfg.environments, &cfg.hooks)
	require.NoError(t, err)

	godump.Fdump(t.Output(), cfg)

	assert.Equal(t, []string{"/bin/bash", "-c", "echo hello world"}, cfg.exec.Cmdline)
	assert.Equal(t, "/root", cfg.exec.Cwd)
	assert.Equal(t, StopMethodKill, cfg.stop.Method)
	assert.Equal(t, []string{"/bin/bash", "-c", "echo \"stop\""}, cfg.stop.Command)
	assert.Equal(t, "SIGTERM", string(cfg.stop.Signal))
	assert.Equal(t, uint(10), cfg.stop.Timeout)
}
