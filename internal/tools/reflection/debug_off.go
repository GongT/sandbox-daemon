//go:build !(debug && reflect)

package reflection

func debug(format string, v ...interface{}) {
	// do nothing
}
