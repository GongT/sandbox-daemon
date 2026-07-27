package config

import (
	"maps"
	"reflect"
	"slices"

	"github.com/tidwall/gjson"
	"gitlab.com/tozd/go/errors"
)

type configJsonContext struct {
	ConfigFillContext

	text string
	json gjson.Result
}

func NewJsonContext(content string) (*configJsonContext, error) {
	result := gjson.Parse(content)
	return &configJsonContext{
		text: content,

		json: result,
	}, nil
}

func (ctx *configJsonContext) HasValue(tagPath ConfigPath) (bool, error) {
	path := tagPath.JoinWithDot("")
	result := ctx.json.Get(path)
	return result.Exists(), nil
}

func (ctx *configJsonContext) GetArraySize(tagPath ConfigPath) (int, error) {
	path := tagPath.JoinWithDot("")
	result := ctx.json.Get(path)
	if !result.Exists() {
		return 0, errors.Errorf("错误调用GetArraySize: 路径不存在")
	}
	if !result.IsArray() {
		return 0, errors.Errorf("不是数组类型")
	}
	return len(result.Array()), nil
}

func (ctx *configJsonContext) GetObjectKeys(tagPath ConfigPath) ([]string, error) {
	path := tagPath.JoinWithDot("")
	result := ctx.json.Get(path)
	if !result.Exists() {
		return nil, errors.Errorf("错误调用GetObjectKeys: 路径不存在")
	}
	if !result.IsObject() {
		return nil, errors.Errorf("不是对象类型")
	}

	keys := slices.Collect(maps.Keys(result.Map()))
	return keys, nil
}

func (ctx *configJsonContext) GetValue(t reflect.Type, tagPath ConfigPath) (string, error) {
	path := tagPath.JoinWithDot("")
	result := ctx.json.Get(path)
	if !result.Exists() {
		return "", errors.Errorf("错误调用GetValue: 路径不存在")
	}
	return result.String(), nil
}
