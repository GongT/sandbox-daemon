package io_streams

import (
	"testing"

	test_handler "github.com/gongt/sandbox-daemon/internal/streams/handlers/test"
	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/stretchr/testify/require"
)

func TestForwarderInherit(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	mux := NewMultiplexer()

	cfg, _ := ioconfig.NewDestination("inherit://stdout")
	closer, err := mux.AddDestination(cfg)
	require.NoError(t, err)
	defer closer.Close()

	cfg, _ = ioconfig.NewDestination("inherit://8")
	closer, err = mux.AddDestination(cfg)
	require.Error(t, err)
	require.Nil(t, closer)
}

func TestForwarder1(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	mux := NewMultiplexer()

	cfg, _ := ioconfig.NewDestination("test://1")
	t1 := test_handler.NewTestWrapSpec(t, cfg)
	closer, err := mux.AddDestination(t1)
	require.NoError(t, err)
	defer closer.Close()

	cfg, _ = ioconfig.NewDestination("test://2")
	t2 := test_handler.NewTestWrapSpec(t, cfg)
	closer, err = mux.AddDestination(t2)
	require.NoError(t, err)
	defer closer.Close()

	cfg, _ = ioconfig.NewDestination("test://3")
	t3 := test_handler.NewTestWrapSpec(t, cfg)
	closer, err = mux.AddDestination(t3)
	require.NoError(t, err)
	defer closer.Close()
}
