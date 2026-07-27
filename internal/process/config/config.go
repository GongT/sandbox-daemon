package config

type HookConfig struct {
	// 命令行
	cmdline []string `config:"cmdline"`
	// 工作目录
	cwd string `config:"cwd"`
	// 是否使用主程序的环境变量，默认是true
	env bool `config:"env"`
}

type HookRunningConfig struct {
	HookConfig

	// 是否使用主程序的namespace和环境变量，默认是true
	namespace bool `config:"namespace"`
}

type hooksConfig struct {
	// 主程序已经在运行之后执行的命令
	running *HookRunningConfig `config:"hooks.running"`
	// 主程序退出之后执行的命令
	stopped *HookConfig `config:"hooks.stopped"`
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
