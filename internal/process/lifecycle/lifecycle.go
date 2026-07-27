package lifecycle

import "github.com/gongt/sandbox-daemon/internal/process/instance"

type ProcessGroupLifecycle struct {
	group *instance.ProcessGroup
	config 
}

func NewProcessGroupWithLifecycle() *ProcessGroupLifecycle {
	return &ProcessGroupLifecycle{
		group: instance.NewProcessGroup(),
	}
}
