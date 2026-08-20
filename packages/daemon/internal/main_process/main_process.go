package main_process

import (
	"io"

	"github.com/pkg/errors"

	"github.com/gongt/sandbox-daemon/packages/process/process_group"
)

var _ io.Closer = (*MainProcess)(nil)

type MainProcess struct {
	pg *process_group.ProcessGroup

	exitCode int
}

func New() *MainProcess {
	return &MainProcess{
		exitCode: -1,
	}
}

func (m *MainProcess) Start() error {
	if m.pg != nil {
		return errors.New("主进程已启动，无法重复启动")
	}

	m.exitCode = -1
	m.pg = process_group.New()

	return nil
}

func (m *MainProcess) Close() error {
	if m == nil || m.pg == nil {
		return nil
	}
	err := m.pg.StopAll()
	m.GetExitCode()
	m.pg = nil
	return err
}

func (m *MainProcess) GetExitCode() (int, error) {
	if m.exitCode >= 0 {
		return m.exitCode, nil
	}
	if m.pg != nil {
		stat, err := m.pg.GetLeader().GetResult()
		if err != nil {
			return -1, err
		}
		c := stat.ExitCode()
		if c < 0 {
			return c, errors.New("主进程退出码无效")
		}
		m.exitCode = c
		return c, nil
	}
	return -1, errors.New("主进程未启动，无法获取退出码")
}
