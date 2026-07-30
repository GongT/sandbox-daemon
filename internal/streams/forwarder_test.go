package io_streams

import (
	"log"
	"testing"

	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
	"github.com/stretchr/testify/require"
)

func TestForwarder(t *testing.T) {
	log.SetOutput(t.Output())

	mux := NewMultiplexer()

	cfg, _ := ioconfig.NewDestination("inherit://stdout")
	closer, err := mux.AddDestination(*cfg)
	require.NoError(t, err)
	defer closer.Close()

	cfg, _ = ioconfig.NewDestination("inherit://8")
	closer, err = mux.AddDestination(*cfg)
	require.Error(t, err)
	require.Nil(t, closer)

}
