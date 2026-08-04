package any_builder

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	log.SetOutput(t.Output())

	b := New()
	b.PrintLine("Hello, The truth is ", 42, '!')

	assert.Equal(t, "Hello, The truth is 42!\n", b.String())
}

func TestIndent(t *testing.T) {
	log.SetOutput(t.Output())

	b := New()
	b.PrintLine("hello {")
	b.Indent()
	b.Print("n")
	b.Print("e")
	b.Print("w")
	b.PrintLine(" world!")
	b.UnIndent()
	b.Print("}")
	b.EnsureNewline()

	assert.Equal(t, "hello {\n\tnew world!\n}\n", b.String())
}
