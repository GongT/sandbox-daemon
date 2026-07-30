package daemon

import (
	"log"
	"maps"
	"slices"

	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
)

type StdioConfig struct {
	Stdin  bool                       `config:"stdin"`
	Stdout []ioconfig.DestinationSpec `config:"stdout"`
	Stderr []ioconfig.DestinationSpec `config:"stderr"`
	All    []ioconfig.DestinationSpec `config:"all"`
}

func (s *StdioConfig) Validate() error {
	if len(s.Stdout) == 0 {
		dest, err := ioconfig.NewDestination("inherit://stdout")
		if err != nil {
			panic(err)
		}
		s.Stdout = append(s.Stdout, *dest)
	}

	if len(s.Stderr) == 0 {
		dest, err := ioconfig.NewDestination("inherit://stderr")
		if err != nil {
			panic(err)
		}
		s.Stderr = append(s.Stderr, *dest)
	}

	s.unique()

	return nil
}

func (s *StdioConfig) unique() {
	stdout := make(map[string]ioconfig.DestinationSpec)
	for _, dest := range s.Stdout {
		key := dest.Raw()
		if _, ok := stdout[key]; !ok {
			stdout[key] = dest
		} else {
			log.Printf("忽略重复写入目标'%s'", key)
		}
	}

	stderr := make(map[string]ioconfig.DestinationSpec)
	for _, dest := range s.Stderr {
		key := dest.Raw()
		if _, ok := stderr[key]; !ok {
			stderr[key] = dest
		} else {
			log.Printf("忽略重复写入目标'%s'", key)
		}
	}

	all := make(map[string]ioconfig.DestinationSpec)
	for _, dest := range s.All {
		key := dest.Raw()
		if _, ok := all[key]; !ok {
			all[key] = dest
		} else {
			log.Printf("忽略重复写入目标'%s'", key)
		}
	}

	for _, dest := range s.Stdout {
		key := dest.Raw()
		if _, ok := stderr[key]; ok {
			log.Printf("写入目标'%s'同时出现在 stdout 和 stderr", key)
			delete(stderr, key)
			all[key] = dest
		}
	}
	for _, dest := range s.Stderr {
		key := dest.Raw()
		if _, ok := stdout[key]; ok {
			log.Printf("写入目标'%s'同时出现在 stdout 和 stderr", key)
			delete(stdout, key)
			all[key] = dest
		}
	}

	s.Stdout = slices.Collect(maps.Values(stdout))
	s.Stderr = slices.Collect(maps.Values(stderr))
	s.All = slices.Collect(maps.Values(all))
}

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
	case "passthrough", "true":
		// all false
	default:
		cfg.IsStop = true
	}
}
