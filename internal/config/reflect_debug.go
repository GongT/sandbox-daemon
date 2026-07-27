package config

import (
	"log"
	"reflect"
)

type loggingContext struct {
	ConfigFillContext

	oCtx ConfigFillContext
}

func (l *loggingContext) HasValue(tagPath ConfigPath) (bool, error) {
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

func (l *loggingContext) GetArraySize(tagPath ConfigPath) (size int, err error) {
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

func (l *loggingContext) GetObjectKeys(tagPath ConfigPath) (keys []string, err error) {
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

func (l *loggingContext) GetValue(t reflect.Type, tagPath ConfigPath) (value string, err error) {
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
