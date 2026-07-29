package config

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
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
		} else {
			return fmt.Errorf("未知停止信号: %d", c.Signal)
		}
		if len(c.Command) > 0 {
			return fmt.Errorf("停止命令方式不支持指定命令")
		}
	case StopMethodCommand:
		if len(c.Command) == 0 {
			return fmt.Errorf("停止命令不能为空")
		}
		if c.Signal != 0 {
			return fmt.Errorf("停止命令方式不支持指定信号")
		}
	default:
		return fmt.Errorf("未知停止方式: %s, 应为 %s 或 %s", c.Method, StopMethodKill, StopMethodCommand)
	}
	return nil
}

type signalFromName syscall.Signal

func (s *signalFromName) FromString(v string) error {
	sVal := unix.SignalNum(v)
	*s = signalFromName(sVal)
	return nil
}
