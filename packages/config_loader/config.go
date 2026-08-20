package config_loader

import (
	"os"
	"strings"

	"github.com/gongt/sandbox-daemon/packages/config_loader/internal"
	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/context"
	context_json "github.com/gongt/sandbox-daemon/packages/config_loader/internal/context/json"
	context_yaml "github.com/gongt/sandbox-daemon/packages/config_loader/internal/context/yaml"
	"github.com/gongt/sandbox-daemon/packages/logger"
	"github.com/gongt/sandbox-daemon/packages/tools/interfaces"
	"github.com/gongt/sandbox-daemon/packages/tools/reflection/deep_init"
	"gitlab.com/tozd/go/errors"
)

func LoadConfigFile(filename string, inputs ...any) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return errors.WithMessagef(err, "读取配置文件%s失败", filename)
	}

	return LoadConfigContent(string(f), inputs...)
}

func LoadConfigContent(content string, inputs ...any) error {
	content = strings.TrimSpace(content)

	var ctx context.ConfigFillContext
	if strings.HasPrefix(content, "{") {
		logger.DConfigF("配置识别为JSON")
		var err error
		ctx, err = context_json.NewJsonContext(content)
		if err != nil {
			return errors.WithMessage(err, "JSON解析错误")
		}
	} else {
		logger.DConfigF("配置识别为YAML")
		var err error
		ctx, err = context_yaml.NewYamlContext(content)
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

func loadConfigInto(input any, ctx context.ConfigFillContext) error {
	walkedValues := deep_init.DeepInitialize(input)
	for _, ptr := range walkedValues {
		if sub, ok := ptr.(interfaces.Initializer); ok {
			sub.Initialize()
		}
	}

	err := internal.WalkStruct(input, ctx)

	if err != nil {
		return err
	}

	for _, ptr := range walkedValues {
		if sub, ok := ptr.(interfaces.Validator); ok {
			err := sub.Validate()
			if err != nil {
				return errors.WithMessage(err, "读取配置文件后验证失败")
			}
		}
	}
	return nil
}
