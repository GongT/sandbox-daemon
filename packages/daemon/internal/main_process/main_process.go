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
		return errors.New("主进程已启动，无法重复启动")
	}
	m.pg = process_group.New()

	return nil
}

func (m *mainProcess) Stop() error {
	if m.pg == nil {
		return errors.New("主进程未启动，无法停止")
	}
	err := m.pg.StopAll()
	m.pg = nil
	return err
}
