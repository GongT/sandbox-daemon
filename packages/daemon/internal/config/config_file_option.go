package daemon_config

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

const DEFAULT_PATH = "/etc/sandbox-daemon.yaml"

type ConfigFileOption struct {
	input    string
	resolved string
}

func (c *ConfigFileOption) GetFilePath() string {
	return c.resolved
}

func (c *ConfigFileOption) IsEmpty() bool {
	return c.input == ""
}

func (c ConfigFileOption) String() string {
	return fmt.Sprintf("{input: %s, resolved: %s}", c.input, c.resolved)
}

func (c *ConfigFileOption) UnmarshalFlag(value string) error {
	c.input = value

	if c.input == "" {
		s, err := os.Stat(DEFAULT_PATH)
		if os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return errors.WithMessagef(err, "无法访问默认配置文件路径: %s", DEFAULT_PATH)
		} else if s.IsDir() {
			return errors.Errorf("默认配置文件路径是一个目录: %s", DEFAULT_PATH)
		}
		c.input = DEFAULT_PATH
	}

	if path.IsAbs(c.input) {
		c.resolved = c.input
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.WithMessage(err, "无法获取当前工作目录")
		}
		c.resolved = path.Join(cwd, c.input)
	}

	return nil
}

type ConfigFileOptionOptional struct {
	ConfigFileOption
}

func (c *ConfigFileOptionOptional) UnmarshalFlag(value string) error {
	c.input = value

	if c.input == "" {
		return nil
	}

	if path.IsAbs(c.input) {
		c.resolved = c.input
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.WithMessage(err, "无法获取当前工作目录")
		}
		c.resolved = path.Join(cwd, c.input)
	}

	return nil
}
