//go:build !release

package print

import "github.com/gongt/sandbox-daemon/packages/logger/internal/tags"

func DebugLogF(t tags.DebugTag, fmt string, v ...any) {
	ReleaseLogF(t, fmt, v...)
}

func DebugLog(t tags.DebugTag, v ...any) {
	ReleaseLog(t, v...)
}
