package config

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	internalconfig "github.com/gongt/sandbox-daemon/internal/config"
)

func TestProcessConfigLoad(t *testing.T) {
	content := []byte(`exec:
  cmdline:
    - /bin/sh
    - -c
    - echo hello
  cwd: /tmp
  stop:
    method: command
    signal: SIGTERM
    command:
      - /bin/true
    timeout: 3
`)

	cfg := New()
	if unused, err := internalconfig.LoadConfigObject(content, cfg); err != nil || len(unused) > 0 {
		if err != nil {
			t.Fatal(err)
		}
		if len(unused) > 0 {
			t.Fatalf("未知配置项: %v", unused)
		}
	}

	spew.Dump(cfg)

	require.Equal(t, []string{"/bin/sh", "-c", "echo hello"}, cfg.cmdline)
	require.Equal(t, "/tmp", cfg.cwd)
	require.Equal(t, StopMethodCommand, cfg.stop.method)
	require.Equal(t, []string{"/bin/true"}, cfg.stop.command)
	require.Equal(t, "SIGTERM", string(cfg.stop.signal))
	require.Equal(t, uint(3), cfg.stop.timeout)

	fmt.Printf("应有错误: %s\n", cfg.Validate())
}
