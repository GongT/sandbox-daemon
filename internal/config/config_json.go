package config

import (
	"maps"
	"reflect"
	"slices"

	"github.com/gongt/sandbox-daemon/internal/tools/i18n/type_name"
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

func (ctx *configJsonContext) GetValue(t reflect.Type, tagPath ConfigPath) (any, error) {
	path := tagPath.JoinWithDot("")
	result := ctx.json.Get(path)
	if !result.Exists() {
		return "", errors.Errorf("错误调用GetValue: 路径不存在")
	}

	switch t.Kind() {
	case reflect.String:
		return result.String(), nil
	case reflect.Bool:
		return result.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return result.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return result.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return result.Float(), nil
	default:
		return nil, errors.Errorf("不支持将%s转换为%s", type_name.TranslateType(reflect.TypeOf(result.Value())), type_name.TranslateType(t))
	}
}
