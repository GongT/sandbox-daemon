package myenv

import (
	"io"
	"log"
	"testing"
)

func RedirectDebugOutput(output io.Writer) {
	if !testing.Testing() {
		panic("RedirectDebugOutput仅用于go test测试")
	}
	log.SetOutput(output)
}

func RedirectDebugTesting(t testing.TB) {
	if !testing.Testing() {
		panic("RedirectDebugTesting仅用于go test测试")
	}
	RedirectDebugOutput(t.Output())
}
