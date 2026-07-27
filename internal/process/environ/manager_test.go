package environ

import (
	"log"
	"reflect"
	"testing"

	"github.com/goforj/godump"
	internalconfig "github.com/gongt/sandbox-daemon/internal/config"
	"github.com/gongt/sandbox-daemon/internal/process/config"
	"github.com/gongt/sandbox-daemon/internal/tools"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentManager(t *testing.T) {
	var err error
	log.SetOutput(t.Output())

	cfg := config.EnvironmentsConfig{}

	content := `
environments:
  add:
    HOME: /tmp
  blacklist:
  - LD_PRELOAD
  whitelist:
  - PATH
  - HOME
`

	err = internalconfig.LoadConfigContent(content, &cfg)
	require.NoError(t, err)
	require.Equal(t, tools.EnvironmentMap{"HOME": "/tmp"}, cfg.Add)

	mgr := New(&cfg)

	mapping := mgr.Snapshot()
	godump.Fdump(t.Output(), mapping)

	require.Equal(t, "string", reflect.TypeOf(mapping["PATH"]).String())
	require.Equal(t, "/tmp", mapping["HOME"])
	require.Equal(t, len(mapping), 2)

	_, exists := mapping["LD_PRELOAD"]
	require.Equal(t, exists, false)
}
