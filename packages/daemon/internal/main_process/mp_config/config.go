package mp_config

import (
	"github.com/gongt/sandbox-daemon/packages/config_loader"
	"github.com/gongt/sandbox-daemon/packages/process/environ"
	"github.com/gongt/sandbox-daemon/packages/process/process_group"
	"github.com/pkg/errors"
)

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
	Running HookRunningConfig `config:"hooks.running"`
	// 主程序退出之后执行的命令
	Stopped HookConfig `config:"hooks.stopped"`
}

type LifecycleConfig struct {
	exec         ExecConfig
	stop         StopConfig
	environments environ.ManagerConfig
	hooks        hooksConfig
	namespace    process_group.NamespaceConfig
}

func NewFromFile(file string) (*LifecycleConfig, error) {
	cfg := &LifecycleConfig{}
	err := cfg.readFrom(file)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func New(text string) (*LifecycleConfig, error) {
	cfg := &LifecycleConfig{}
	err := cfg.loadFrom(text)
	if err != nil {
		return nil, err
	}
	return cfg, nil
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

func (c *LifecycleConfig) readFrom(file string) error {
	err := config_loader.LoadConfigFile(file, &c.exec, &c.stop, &c.environments, &c.hooks, &c.namespace)
	if err != nil {
		return errors.WithMessage(err, "程序配置错误")
	}
	return nil
}

func (c *LifecycleConfig) loadFrom(text string) error {
	err := config_loader.LoadConfigContent(text, &c.exec, &c.stop, &c.environments, &c.hooks, &c.namespace)
	if err != nil {
		return errors.WithMessage(err, "程序配置错误")
	}
	return nil
}

func (c *LifecycleConfig) HasExecSpecified() bool {
	return c.exec.HasSpecified()
}
