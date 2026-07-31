package logger

import (
	"github.com/gongt/sandbox-daemon/packages/logger/internal/tags"
	"github.com/gongt/sandbox-daemon/packages/logger/print"
)

func DLogF(t tags.DebugTag, fmt string, v ...any) {
	print.DebugLogF(t, fmt, v...)
}

func DLog(t tags.DebugTag, v ...any) {
	print.DebugLog(t, v...)
}

func LogF(t tags.DebugTag, fmt string, v ...any) {
	print.ReleaseLogF(t, fmt, v...)
}

func Log(t tags.DebugTag, v ...any) {
	print.ReleaseLog(t, v...)
}
