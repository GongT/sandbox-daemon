package config

import (
	"os"
	"strings"

	"github.com/gongt/sandbox-daemon/internal/tools"
	"github.com/gongt/sandbox-daemon/internal/tools/interfaces"
	"github.com/gongt/sandbox-daemon/internal/tools/reflection/deep_init"
	"gitlab.com/tozd/go/errors"
)

func LoadConfigFile(filename string, inputs ...interface{}) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return errors.WithMessagef(err, "读取配置文件%s失败", filename)
	}

	return LoadConfigContent(string(f), inputs...)
}

func LoadConfigContent(content string, inputs ...interface{}) error {
	content = strings.TrimSpace(content)

	var ctx ConfigFillContext
	if strings.HasPrefix(content, "{") {
		tools.DebugLog("配置识别为JSON")
		var err error
		ctx, err = NewJsonContext(content)
		if err != nil {
			return errors.WithMessage(err, "JSON解析错误")
		}
	} else {
		tools.DebugLog("配置识别为YAML")
		var err error
		ctx, err = NewYamlContext(content)
		if err != nil {
			return errors.WithMessage(err, "YAML解析错误")
		}
	}

	for _, input := range inputs {
		if err := loadConfigInto(input, ctx); err != nil {
			return err
		}
	}
	return nil
}

func loadConfigInto(input any, ctx ConfigFillContext) error {
	walkedValues := deep_init.DeepInitialize(input)
	for _, ptr := range walkedValues {
		if sub, ok := ptr.(interfaces.Initializer); ok {
			sub.Initialize()
		}
	}

	err := WalkStruct(input, ctx)

	if err != nil {
		return err
	}

	for _, ptr := range walkedValues {
		if sub, ok := ptr.(interfaces.Validator); ok {
			err := sub.Validate()
			if err != nil {
				return errors.WithMessagef(err, "读取配置文件后验证失败")
			}
		}
	}
	return nil
}
