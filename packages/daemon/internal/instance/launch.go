package instance

import (
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process/mp_config"
	"github.com/gongt/sandbox-daemon/packages/tools/signals"
)

// 启动主进程
func (d *D) LaunchMainProcess(config *mp_config.LifecycleConfig) error {
	err := d.mp.Start()
	if err != nil {
		return err
	}

	if d.config.ForwardCode {
		code, _ := d.mp.GetExitCode()
		signals.AppQuit.Set(code)
	}

	return nil
}

// 退出主进程
func (d *D) QuitMainProcess() error {
	return d.mp.Close()
}
