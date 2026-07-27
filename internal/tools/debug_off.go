//go:build !debug

package tools

const IsDebug = false

func DebugLog(format string, v ...interface{}) {
	// do nothing
}

func DebugLogOnce(format string, v ...interface{}) {
	// do nothing
}
