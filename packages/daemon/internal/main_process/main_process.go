package main_process

import (
	"errors"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	"github.com/gongt/sandbox-daemon/packages/process/process_group"
)

var _ internal.DaemonComponent = (*mainProcess)(nil)

type mainProcess struct {
	pg *process_group.ProcessGroup
}

func New() *mainProcess {
	return &mainProcess{}
}

func (m *mainProcess) Start() error {
	if m.pg != nil {
		return errors.New("main process already started")
	}
	m.pg = process_group.New()

	return nil
}

func (m *mainProcess) Stop() error {
	err := m.pg.StopAll()
	m.pg = nil
	return err
}
