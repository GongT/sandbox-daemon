package environ

import (
	"log"
	"testing"

	"github.com/goforj/godump"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentManager(t *testing.T) {
	log.SetOutput(t.Output())

	cfg := ManagerConfig{
		Add: map[string]string{
			"TEST": "VALUE",
		},
		Blacklist: []string{
			"LD_PRELOAD",
		},
		Whitelist: []string{
			"PATH",
			"HOME",
		},
		Requires: []string{
			"MISSING",
		},
	}

	mgr := NewManager(&cfg)
	mgr.initial = map[string]string{
		"PATH":       "/usr/bin",
		"HOME":       "/tmp",
		"LD_PRELOAD": "lib.so",
	}

	mapping, err := mgr.Snapshot()
	godump.Fdump(t.Output(), mapping)

	require.Equal(t, Map{
		"PATH": "/usr/bin",
		"HOME": "/tmp",
		"TEST": "VALUE",
	}, mapping)

	_, exists := mapping["LD_PRELOAD"]
	require.Equal(t, exists, false)

	require.Error(t, err)
}
