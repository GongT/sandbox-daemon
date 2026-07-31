package print

import (
	"log"

	"github.com/goforj/godump"
	"github.com/gongt/sandbox-daemon/packages/logger/internal/tags"
)

func ReleaseLogF(t tags.DebugTag, fmt string, v ...any) {
	if !tags.CheckEnabled(t) {
		return
	}
	log.Printf(fmt, v...)
}

func ReleaseLog(t tags.DebugTag, v ...any) {
	if !tags.CheckEnabled(t) {
		return
	}
	log.Print(godumpAll(v))
}

func godumpAll(v []any) []any {
	result := make([]any, len(v))
	for i, x := range v {
		result[i] = godump.DumpStr(x)
	}
	return result
}
