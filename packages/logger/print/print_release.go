//go:build release

package print

import "github.com/gongt/sandbox-daemon/packages/logger/internal/tags"

func DebugLogF(t tags.DebugTag, fmt string, v ...any) {
	// empty
}

func DebugLog(t tags.DebugTag, v ...any) {
	// empty
}
