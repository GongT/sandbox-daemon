package daemon_config

type DaemonConfig struct {
	RuntimeFolder string

	ForwardCode   bool         `config:"daemon.exit_code"`
	SigintAction  SignalAction `config:"daemon.sigint"`
	SigtermAction SignalAction `config:"daemon.sigterm"`
	StdioConfig   StdioConfig  `config:"daemon.stdio"`
}

type SignalAction struct {
	IsIgnore bool
	IsStop   bool
}

func (cfg *SignalAction) FromString(s string) {
	switch s {
	case "ignore", "false":
		cfg.IsIgnore = true
	case "passthrough":
		// noop
	default:
		cfg.IsStop = true
	}
}
