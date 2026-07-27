package config

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/gongt/sandbox-daemon/internal/tools"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func LoadConfigFile(filename string, inputs ...interface{}) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return LoadConfigContent(string(f), inputs...)
}

func LoadConfigContent(content string, inputs ...interface{}) error {
	content = strings.TrimSpace(content)

	var rv = viper.New()
	if strings.HasPrefix(content, "{") {
		tools.DebugLog("配置识别为JSON")
		rv.SetConfigType("json")
	} else {
		tools.DebugLog("配置识别为YAML")
		rv.SetConfigType("yaml")
	}
	err := rv.ReadConfig(bytes.NewBuffer([]byte(content)))
	if err != nil {
		return errors.WithStack(fmt.Errorf("读取配置失败: %s", err.Error()))
	}

	for _, input := range inputs {
		if err := loadConfigInto(input, rv); err != nil {
			return err
		}
	}
	return nil
}

type viperContext struct {
	ConfigFillContext

	rv *viper.Viper
}

func (v *viperContext) HasValue(tagPath []string) (bool, error) {
	key := strings.Join(tagPath, ".")
	r := v.rv.IsSet(key)

	for _, tag := range tagPath {
		if containsUppercase(tag) {
			return false, errors.WithStack(fmt.Errorf("配置tag不能包含大写字母: %q", tag))
		}
	}

	return r, nil
}

func containsUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' { // key不能有英文大写字母
			return true
		}
	}
	return false
}

func (v *viperContext) GetArraySize(tagPath []string) (int, error) {
	key := strings.Join(tagPath, ".")
	if !v.rv.IsSet(key) {
		return 0, errors.WithStack(fmt.Errorf("错误调用GetArraySize: 路径不存在"))
	}
	data := v.rv.Get(key)
	typ := reflect.TypeOf(data)
	if typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
		return 1, nil
	}
	return reflect.ValueOf(data).Len(), nil
}

func (v *viperContext) GetObjectKeys(tagPath []string) ([]string, error) {
	key := strings.Join(tagPath, ".")
	if !v.rv.IsSet(key) {
		return nil, errors.WithStack(fmt.Errorf("错误调用GetObjectKeys: 路径不存在"))
	}
	data := v.rv.GetStringMapString(key)
	return slices.Collect(maps.Keys(data)), nil
}

func (v *viperContext) GetValue(t reflect.Type, tagPath []string) (string, error) {
	key := strings.Join(tagPath, ".")
	if !v.rv.IsSet(key) {
		return "", errors.WithStack(fmt.Errorf("错误调用GetValue: 路径不存在"))
	}
	return v.rv.GetString(key), nil
}

func loadConfigInto(input interface{}, rv *viper.Viper) error {
	ctx := &viperContext{
		rv: rv,
	}

	if err := WalkStruct(input, ctx); err != nil {
		return err
	}

	return nil
}
