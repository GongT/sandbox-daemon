//go:build debug

package tools

import (
	"fmt"
	"log"
	"runtime"
)

const IsDebug = true

func DebugLog(format string, v ...interface{}) {
	log.Printf(format, v...)
}

var debugLogOnceMap = make(map[string]bool)

func DebugLogOnce(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	key := fmt.Sprintf("%s:%d", file, line)
	if !debugLogOnceMap[key] {
		debugLogOnceMap[key] = true
		log.Printf(format, v...)
	}
}
