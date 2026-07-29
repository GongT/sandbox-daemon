//go:build debug && reflect

package deep_init

import (
	"log"
)

func debug(format string, v ...interface{}) {
	log.Printf(format, v...)
}
