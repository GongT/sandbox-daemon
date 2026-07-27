package config

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/goforj/godump"
	"github.com/stretchr/testify/require"

	assets "github.com/gongt/sandbox-daemon"
	internalconfig "github.com/gongt/sandbox-daemon/internal/config"
)

func TestProcessConfigLoad(t *testing.T) {
	content := []byte(assets.ExampleConfigYAML)

	cfg := LifecycleConfig{}
	if unused, err := internalconfig.LoadConfigObject(content, &cfg.exec, &cfg.stop, &cfg.environments, &cfg.hooks); err != nil || len(unused) > 0 {
		if err != nil {
			t.Fatal(err)
		}
		if len(unused) > 0 {
			for _, name := range unused {
				t.Errorf("未知配置项: %v", name)
			}
		}
	}

	godump.Fdump(t.Output(), cfg)

	require.Equal(t, []string{"/bin/bash", "-c", "echo hello world"}, cfg.exec.cmdline)
	require.Equal(t, "/root", cfg.exec.cwd)
	require.Equal(t, StopMethodCommand, cfg.stop.method)
	require.Equal(t, []string{"/bin/bash", "-c", "echo \"stop\""}, cfg.stop.command)
	require.Equal(t, "SIGTERM", string(cfg.stop.signal))
	require.Equal(t, uint(10), cfg.stop.timeout)

	fmt.Printf("应有错误: %s\n", cfg.Validate())
}
