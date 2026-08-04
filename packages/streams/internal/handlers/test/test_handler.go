package test_handler

import (
	"log"
	"strings"
	"testing"

	"github.com/gongt/sandbox-daemon/packages/streams/internal/ioconfig"
	"github.com/gongt/sandbox-daemon/packages/streams/internal/linetype"
)

var _ ioconfig.LineDestination = (*testForwarder)(nil)

type testWrapSpec struct {
	ioconfig.DestinationSpec
	T      *testing.T
	Output strings.Builder
}

func NewTestWrapSpec(t *testing.T, cfg ioconfig.DestinationSpec) *testWrapSpec {
	return &testWrapSpec{
		DestinationSpec: cfg,
		T:               t,
		Output:          strings.Builder{},
	}
}

type testForwarder struct {
	t *testing.T
	o *strings.Builder
}

func NewTestForwarder() *testForwarder {
	return &testForwarder{}
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
	ws := target.(*testWrapSpec)
	self.t = ws.T
	self.o = &ws.Output
	return nil
}
