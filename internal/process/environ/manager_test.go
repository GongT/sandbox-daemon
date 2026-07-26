package environ

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	internalconfig "github.com/gongt/sandbox-daemon/internal/config"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentManager(t *testing.T) {
	cfg, err := NewEnvironmentManager()
	require.NoError(t, err)

	content := []byte(`
environments:
  add:
    - HOME=/tmp
  blacklist:
    - LD_PRELOAD
  whitelist:
    - PATH
    - HOME
`)

	if unused, err := internalconfig.LoadConfigObject(content, cfg); err != nil || len(unused) > 0 {
		if err != nil {
			t.Fatal(err)
		}
		if len(unused) > 0 {
			t.Fatalf("未知配置项: %v", unused)
		}
	}

	require.NoError(t, err)

	mapping := cfg.Snapshot()
	spew.Dump(mapping)

	require.Equal(t, "string", reflect.TypeOf(mapping["PATH"]).String())
	require.Equal(t, "/tmp", mapping["HOME"])
	require.Equal(t, len(mapping), 2)

	_, exists := mapping["LD_PRELOAD"]
	require.Equal(t, exists, false)
}
