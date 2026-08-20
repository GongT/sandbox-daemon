package daemon

import (
	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	"github.com/gongt/sandbox-daemon/packages/myenv"
)

type PipeCommand struct {
	internal.WithSessionCommand

	Readonly bool `long:"readonly" description:"不要接入stdin" default:"false"`
}

func (config *PipeCommand) Run(runtime *myenv.GlobalOptions) error {
	return nil
}
