package test_handler

import (
	"log"

	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
	"github.com/gongt/sandbox-daemon/internal/streams/linetype"
)

var _ ioconfig.LineDestination = (*testForwarder)(nil)

type testForwarder struct {
}

func (self *testForwarder) Execute(ch <-chan linetype.LineData) error {
	for line := range ch {
		log.Println(line.Line)
	}

	return nil
}

func (self *testForwarder) Destroy() error {
	return nil
}

func (self *testForwarder) Initialize(target ioconfig.DestinationSpec) error {
	return nil
}
