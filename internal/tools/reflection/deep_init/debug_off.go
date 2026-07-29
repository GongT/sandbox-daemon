//go:build !(debug && reflect)

package deep_init

func debug(format string, v ...interface{}) {
	// do nothing
}
