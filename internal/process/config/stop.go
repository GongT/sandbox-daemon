package config

import (
	"fmt"

	"github.com/gongt/sandbox-daemon/internal/process/signal"
)

type StopConfig struct {
	method  StopMethod        `config:"stop.method"`
	command []string          `config:"stop.command"`
	signal  signal.SignalName `config:"stop.signal"`
	timeout uint              `config:"stop.timeout"`
}

func (c *StopConfig) Validate() error {
	switch c.method {
	case StopMethodKill:
		if c.signal == "" {
			c.signal = "SIGTERM"
		} else {
			if !c.signal.IsValid() {
				return fmt.Errorf("未知停止信号: %s", c.signal)
			}
		}
		if len(c.command) > 0 {
			return fmt.Errorf("停止命令方式不支持指定命令")
		}
	case StopMethodCommand:
		if len(c.command) == 0 {
			return fmt.Errorf("停止命令不能为空")
		}
		if c.signal != "" {
			return fmt.Errorf("停止命令方式不支持指定信号")
		}
	default:
		return fmt.Errorf("未知停止方式: %s, 应为 %s 或 %s", c.method, StopMethodKill, StopMethodCommand)
	}
	return nil
}
