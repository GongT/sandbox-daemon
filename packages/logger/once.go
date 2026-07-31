package logger

import "github.com/gongt/sandbox-daemon/packages/logger/internal/once"

func DebugLogOnce(format string, v ...interface{}) {
	once.DebugLogOnce(2, format, v...)
}

func AlertOnce(format string, v ...interface{}) {
	once.AlertOnce(2, format, v...)
}
