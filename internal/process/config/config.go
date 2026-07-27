package config

type HookConfig struct {
	// 命令行
	Cmdline []string `config:"cmdline"`
	// 工作目录
	Cwd string `config:"cwd"`
	// 是否使用主程序的环境变量，默认是true
	Env bool `config:"env"`
}

type HookRunningConfig struct {
	HookConfig

	// 是否使用主程序的Namespace和环境变量，默认是true
	Namespace bool `config:"namespace"`
}

type hooksConfig struct {
	// 主程序已经在运行之后执行的命令
	Running *HookRunningConfig `config:"hooks.running"`
	// 主程序退出之后执行的命令
	Stopped *HookConfig `config:"hooks.stopped"`
}

type LifecycleConfig struct {
	exec         ExecConfig
	stop         StopConfig
	environments EnvironmentsConfig
	hooks        hooksConfig
}

func (c *LifecycleConfig) Validate() error {
	if err := c.exec.Validate(); err != nil {
		return err
	}
	if err := c.stop.Validate(); err != nil {
		return err
	}
	return nil
}
