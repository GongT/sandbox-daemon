package config

import (
	_ "embed"
	"testing"

	"github.com/goforj/godump"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	assets "github.com/gongt/sandbox-daemon"
	"github.com/gongt/sandbox-daemon/packages/config_loader"
	"github.com/gongt/sandbox-daemon/packages/myenv"
)

func TestProcessConfigLoad(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	cfg := LifecycleConfig{}
	err := config_loader.LoadConfigContent(assets.ExampleConfigYAML, &cfg.exec, &cfg.stop, &cfg.environments, &cfg.hooks)
	require.Error(t, err, "不支持指定")

	godump.Fdump(t.Output(), cfg)

	assert.Equal(t, []string{"/bin/bash", "-c", "echo hello world"}, cfg.exec.Cmdline)
	assert.Equal(t, "/root", cfg.exec.Cwd)
	assert.Equal(t, StopMethodKill, cfg.stop.Method)
	assert.Equal(t, []string{"/bin/bash", "-c", "echo \"stop\""}, cfg.stop.Command)
	assert.Equal(t, unix.SIGTERM, unix.Signal(cfg.stop.Signal))
	assert.Equal(t, uint(10), cfg.stop.Timeout)
}

func TestInvalidSignal(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	yml := `
stop:
  method: kill
  signal: 9999
`

	cfg := LifecycleConfig{}
	err := config_loader.LoadConfigContent(yml, &cfg.exec)
	require.Error(t, err)
}
