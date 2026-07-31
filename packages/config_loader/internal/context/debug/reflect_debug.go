package context_debug

import (
	"log"
	"reflect"

	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/context"
	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/paths"
)

type loggingContext struct {
	context.ConfigFillContext

	oCtx context.ConfigFillContext
}

func NewLoggingContext(oCtx context.ConfigFillContext) *loggingContext {
	return &loggingContext{
		oCtx: oCtx,
	}
}

func (l *loggingContext) HasValue(tagPath paths.ConfigPath) (bool, error) {
	if true {
		log.Printf("[HasValue] %s", tagPath.JoinWithDot(""))
	}
	r, err := l.oCtx.HasValue(tagPath)

	if err != nil {
		log.Printf("           -> %t, error: %.50s", r, err.Error())
	} else {
		log.Printf("           -> %t", r)
	}
	return r, err
}

func (l *loggingContext) GetArraySize(tagPath paths.ConfigPath) (size int, err error) {
	if true {
		log.Printf("[getArraySize]  %s", tagPath.JoinWithDot(""))
	}
	size, err = l.oCtx.GetArraySize(tagPath)
	if err != nil {
		log.Printf("                -> error: %.50s", err.Error())
	} else {
		log.Printf("                -> size: %d", size)
	}
	return size, err
}

func (l *loggingContext) GetObjectKeys(tagPath paths.ConfigPath) (keys []string, err error) {
	if true {
		log.Printf("[getObjectKeys] %s", tagPath.JoinWithDot(""))
	}
	keys, err = l.oCtx.GetObjectKeys(tagPath)
	if err != nil {
		log.Printf("                -> error: %.50s", err.Error())
	} else {
		log.Printf("                -> keys: %v", keys)
	}
	return keys, err
}

func (l *loggingContext) GetValue(t reflect.Type, tagPath paths.ConfigPath) (value any, err error) {
	if true {
		log.Printf("[getValue] %s | %s", tagPath.JoinWithDot(""), t.String())
	}
	value, err = l.oCtx.GetValue(t, tagPath)
	if err != nil {
		log.Printf("           -> error: %.50s", err.Error())
	} else {
		log.Printf("           -> value: %s", value)
	}
	return value, err
}
