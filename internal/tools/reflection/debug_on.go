//go:build debug && reflect

package reflection

import (
	"log"
)

func debug(format string, v ...interface{}) {
	log.Printf(format, v...)
}
