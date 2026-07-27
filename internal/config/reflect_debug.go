package config

import (
	"log"
	"reflect"
	"strings"
)

type loggingContext struct {
	ConfigFillContext

	oCtx ConfigFillContext
}

func (l *loggingContext) HasValue(tagPath []string) (bool, error) {
	if true {
		log.Printf("[HasValue] %s", strings.Join(tagPath, "."))
	}
	r, err := l.oCtx.HasValue(tagPath)

	if err != nil {
		log.Printf("           -> %t, error: %.50s", r, err.Error())
	} else {
		log.Printf("           -> %t", r)
	}
	return r, err
}

func (l *loggingContext) GetArraySize(tagPath []string) (size int, err error) {
	if true {
		log.Printf("[getArraySize]  %s", strings.Join(tagPath, "."))
	}
	size, err = l.oCtx.GetArraySize(tagPath)
	if err != nil {
		log.Printf("                -> error: %.50s", err.Error())
	} else {
		log.Printf("                -> size: %d", size)
	}
	return size, err
}

func (l *loggingContext) GetObjectKeys(tagPath []string) (keys []string, err error) {
	if true {
		log.Printf("[getObjectKeys] %s", strings.Join(tagPath, "."))
	}
	keys, err = l.oCtx.GetObjectKeys(tagPath)
	if err != nil {
		log.Printf("                -> error: %.50s", err.Error())
	} else {
		log.Printf("                -> keys: %v", keys)
	}
	return keys, err
}

func (l *loggingContext) GetValue(t reflect.Type, tagPath []string) (value string, err error) {
	if true {
		log.Printf("[getValue] %s | %s", strings.Join(tagPath, "."), t.String())
	}
	value, err = l.oCtx.GetValue(t, tagPath)
	if err != nil {
		log.Printf("           -> error: %.50s", err.Error())
	} else {
		log.Printf("           -> value: %s", value)
	}
	return value, err
}

func (l *loggingContext) ConvertNonScalar(get_value string, t reflect.Type) (interface{}, error) {
	if true {
		log.Printf("[ConvertNonScalar] %s | %s", get_value, t.String())
	}
	r, err := l.oCtx.ConvertNonScalar(get_value, t)
	if err != nil {
		log.Printf("           -> error: %.50s", err.Error())
	} else {
		log.Printf("           -> value: %v", r)
	}

	return r, err
}
