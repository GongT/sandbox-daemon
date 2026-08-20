package mp_config

import (
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type StopMethod string

const (
	StopMethodKill    StopMethod = "kill"
	StopMethodCommand StopMethod = "command"
)

type StopConfig struct {
	Method  StopMethod     `config:"stop.method"`
	Command []string       `config:"stop.command"`
	Signal  signalFromName `config:"stop.signal"`
	Timeout uint           `config:"stop.timeout"`
}

func (c *StopConfig) Validate() error {
	switch c.Method {
	case StopMethodKill:
		if c.Signal == 0 {
			c.Signal = signalFromName(unix.SIGTERM)
		}
		if len(c.Command) > 0 {
			return errors.Errorf("停止方式为 kill 时不支持指定 command")
		}
	case StopMethodCommand:
		if len(c.Command) == 0 {
			return errors.Errorf("停止方式为 command 时 command 不能为空")
		}
		if c.Signal != 0 {
			return errors.Errorf("停止方式为 command 时不支持指定 signal")
		}
	default:
		return errors.Errorf("未知停止方式: %s, 应为 %s 或 %s", c.Method, StopMethodKill, StopMethodCommand)
	}
	return nil
}

type signalFromName syscall.Signal

func (s *signalFromName) FromString(v string) error {
	sVal := unix.SignalNum(v)
	*s = signalFromName(sVal)
	if sVal == 0 {
		return errors.Errorf("未知停止信号: %s", v)
	}
	return nil
}

func (s signalFromName) Validate() error {
	if unix.SignalName(syscall.Signal(s)) == "" {
		return errors.Errorf("未知停止信号: %d", s)
	}
	return nil
}
