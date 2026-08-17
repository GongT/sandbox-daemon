package io_streams

import "github.com/gongt/sandbox-daemon/packages/streams/internal/ioconfig"

type DestinationSpec = ioconfig.DestinationSpec

func NewDestination(spec string) (DestinationSpec, error) {
	return ioconfig.NewDestination(spec)
}
