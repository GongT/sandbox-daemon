package mp_config

type ExecConfig struct {
	Cmdline []string `config:"exec.cmdline"`
	Cwd     string   `config:"exec.cwd"`
}

func (cfg ExecConfig) Validate() error {
	return nil
}

func (cfg ExecConfig) HasSpecified() bool {
	return len(cfg.Cmdline) > 0
}
