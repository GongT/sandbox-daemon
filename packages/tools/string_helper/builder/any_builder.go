package any_builder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gongt/sandbox-daemon/packages/tools/interfaces"
	"github.com/gongt/sandbox-daemon/packages/tools/reflection"
)

type ErrorAction int

const (
	Ignore ErrorAction = iota
	Collect
	Panic
)

type VendorAction int

const (
	SerializeDump VendorAction = iota
	SerializeJson
)

type serializer func(v any) ([]byte, error)

type bWriter struct {
	parent *AnyBuilder
}

func (b *bWriter) Write(p []byte) (n int, err error) {
	return b.parent.writeString(string(p))
}

type AnyBuilder struct {
	builder strings.Builder
	W       io.Writer

	OnError    ErrorAction
	Serializer serializer
	Default    VendorAction
	errors     []error

	indent      []byte
	isLineStart bool
}

func New() *AnyBuilder {
	b := &AnyBuilder{
		OnError: Panic,
	}
	b.W = &bWriter{
		parent: b,
	}
	return b
}

func (b *AnyBuilder) onError(err error) {
	switch b.OnError {
	case Ignore:
		return
	case Collect:
		b.errors = append(b.errors, err)
	case Panic:
		panic(err)
	}
}

func (b *AnyBuilder) GetError() error {
	if len(b.errors) == 0 {
		return nil
	}
	err := errors.Join(b.errors...)
	b.errors = nil
	return err
}

func (b *AnyBuilder) Indent() {
	b.indent = append(b.indent, '\t')
}

func (b *AnyBuilder) UnIndent() {
	if len(b.indent) > 0 {
		b.indent = b.indent[:len(b.indent)-1]
	}
}

func (b *AnyBuilder) Print(args ...any) {
	for _, element := range args {
		ptr := reflection.IndirectAny(element)

		switch v := ptr.(type) {
		case []byte:
			b.writeString(string(v))
		case string:
			b.writeString(v)
		case byte:
			b.writeString(string(v))
		case rune:
			b.writeString(string(v))
		case strings.Builder:
			b.writeString(v.String())
		case bytes.Buffer:
			b.writeString(v.String())
		default:
			if b.Serializer != nil {
				data, err := b.Serializer(v)
				if err != nil {
					b.onError(err)
				} else {
					b.writeString(string(data))
				}
				continue
			}

			switch b.Default {
			case SerializeDump:
				switch v := v.(type) {
				case fmt.Stringer:
					b.writeString(v.String())
				case interfaces.StringerE:
					if str, err := v.String(); err != nil {
						b.onError(err)
					} else {
						b.writeString(str)
					}
				default:
					b.writeString(fmt.Sprintf("%#v", v))
				}
			case SerializeJson:
				data, err := json.Marshal(v)
				if err != nil {
					b.onError(err)
				} else {
					b.writeString(string(data))
				}
			default:
				b.writeString(fmt.Sprintf("%#v", v))
			}
		}
	}
}

func (b *AnyBuilder) PrintLine(s ...any) {
	b.Print(s...)
	b.writeString("\n")
}

func (b *AnyBuilder) Format(format string, args ...any) {
	b.writeString(fmt.Sprintf(format, args...))
}

func (b *AnyBuilder) FormatLine(format string, args ...any) {
	b.Format(format, args...)
	b.writeString("\n")
}

func (b *AnyBuilder) String() string {
	return b.builder.String()
}

func (b *AnyBuilder) writeString(s string) (int, error) {
	var i int
	for line := range strings.Lines(s) { // line 是带换行符的
		if len(line) == 0 {
			continue
		}
		if b.isLineStart {
			ii, _ := b.builder.Write(b.indent)
			i += ii
			b.isLineStart = false
		}

		ii, _ := b.builder.WriteString(line)
		i += ii
		if line[len(line)-1] == '\n' {
			// 换行符结尾，下一行需要缩进
			b.isLineStart = true
		} else {
			b.isLineStart = false
		}
	}
	return i, nil
}

func (b *AnyBuilder) EnsureNewline() {
	if !b.isLineStart {
		b.writeString("\n")
		b.isLineStart = true
	}
}
