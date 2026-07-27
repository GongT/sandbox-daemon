package config

import (
	"os"
	"strings"

	"github.com/gongt/sandbox-daemon/internal/tools"
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

	var ctx *configYamlContext
	if strings.HasPrefix(content, "{") {
		return errors.Errorf("配置识别为JSON, 但JSON格式暂不支持")
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

func loadConfigInto(input interface{}, ctx *configYamlContext) error {
	return WalkStruct(input, ctx)
}
