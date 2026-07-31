package once

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

var onceMap sync.Map

func firstTime(file string, line int) bool {
	key := fmt.Sprintf("%s:%d", file, line)

	_, loaded := onceMap.LoadOrStore(key, true)
	return !loaded
}

// 仅一次，输出一条信息到log.Default()
// 根据调用位置的文件名和行号判断是否首次调用
// 如果首次调用，返回true
func AlertOnce(skip int, format string, v ...any) {
	_, file, line, _ := runtime.Caller(skip)
	if firstTime(file, line) {
		log.Printf(format, v...)
	}
}
