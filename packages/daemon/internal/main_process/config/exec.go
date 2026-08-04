package config

type ExecConfig struct {
	Cmdline []string `config:"exec.cmdline"`
	Cwd     string   `config:"exec.cwd"`
}

func (cfg ExecConfig) Validate() error {
	// if len(cfg.Cmdline) == 0 {
	// 	return errors.Errorf("执行命令不能为空")
	// }

	return nil
}
