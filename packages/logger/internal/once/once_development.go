//go:build !release

package once

import (
	"log"
	"runtime"
)

// 仅一次，输出一条信息到log.Default()
func DebugLogOnce(skip int, format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(skip)
	if firstTime(file, line) {
		log.Printf(format, v...)
	}
}
