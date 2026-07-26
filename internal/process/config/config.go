package config

import (
	"fmt"

	"github.com/gongt/sandbox-daemon/internal/process/signal"
)

type StopMethod string

const (
	StopMethodKill    StopMethod = "kill"
	StopMethodCommand StopMethod = "command"
)

type ExecConfig struct {
	cmdline []string   `config:"exec.cmdline"`
	cwd     string     `config:"exec.cwd"`
	stop    StopConfig `config:"exec.stop"`
}

type StopConfig struct {
	method  StopMethod        `config:"method"`
	command []string          `config:"command"`
	signal  signal.SignalName `config:"signal"`
	timeout uint              `config:"timeout"`
}

func New() *ExecConfig {
	cfg := &ExecConfig{
		stop: StopConfig{
			method:  StopMethodCommand,
			signal:  "",
			timeout: 10,
		},
	}
	return cfg
}

func (cfg ExecConfig) Validate() error {
	switch cfg.stop.method {
	case StopMethodKill:
		if cfg.stop.signal == "" {
			cfg.stop.signal = "SIGTERM"
		} else {
			if !cfg.stop.signal.IsValid() {
				return fmt.Errorf("未知停止信号: %s", cfg.stop.signal)
			}
		}
		if len(cfg.stop.command) > 0 {
			return fmt.Errorf("停止命令方式不支持指定命令")
		}
	case StopMethodCommand:
		if len(cfg.stop.command) == 0 {
			return fmt.Errorf("停止命令不能为空")
		}
		if cfg.stop.signal != "" {
			return fmt.Errorf("停止命令方式不支持指定信号")
		}
	default:
		return fmt.Errorf("未知停止方式: %s, 应为 %s 或 %s", cfg.stop.method, StopMethodKill, StopMethodCommand)
	}

	if len(cfg.cmdline) == 0 {
		return fmt.Errorf("执行命令不能为空")
	}

	return nil
}
