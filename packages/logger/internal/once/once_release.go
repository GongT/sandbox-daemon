//go:build release

package once

// release开启时，DebugLogOnce不输出任何日志
func DebugLogOnce(skip int, format string, v ...interface{}) {
	// empty
}
