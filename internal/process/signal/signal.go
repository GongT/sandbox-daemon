package signal

import (
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

type SignalName string

func (s SignalName) Number() unix.Signal {
	name := string(s)
	if len(name) > 3 && strings.ToUpper(name[:3]) != "SIG" {
		name = "SIG" + name
	}
	return unix.SignalNum(name)
}

func (s SignalName) IsValid() bool {
	return s.Number() != 0
}

func New(input string) SignalName {
	if num, err := strconv.ParseUint(input, 10, 32); err == nil {
		return SignalName(unix.SignalName(unix.Signal(num)))
	}
	return SignalName(input)
}
