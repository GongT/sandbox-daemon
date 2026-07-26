package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type exampleConfig struct {
	stringVal  string           `config:"a.value"`
	stringsVal []string         `config:"a.value2"`
	numberVal  int              `config:"a.number"`
	numbersVal []int            `config:"a.number2"`
	boolVal    bool             `config:"a.bool"`
	subConfig  exampleSubConfig `config:"b"`
	rootConfig exampleConfig2

	mapped       map[string]uint32 `config:"mapped"`
	customLoader environmentLoader `config:"env"`
}

type environmentLoader struct {
	values map[string]string
}

func (c *environmentLoader) Unmarshal(node *yaml.Node) error {
	if c.values == nil {
		c.values = map[string]string{}
	}

	if node.Kind == yaml.MappingNode {
		return node.Decode(c.values)
	}

	arr := []string{}

	if err := LoadArray(node, &arr); err != nil {
		return err
	}

	for _, v := range arr {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			c.values[parts[0]] = parts[1]
		} else {
			return fmt.Errorf("无法解析环境变量: %s", v)
		}
	}

	return nil
}

type exampleSubConfig struct {
	stringVal string `config:"value"`
}

type exampleConfig2 struct {
	unknownField int
	x            int `config:"x"`
}

func TestConfig(t *testing.T) {
	exampleYamlContent := `
a:
  value: "hello"
  value2: "world"
  number: "42"
  number2: 
    - 43
    - "44"
  bool: true
  bool: 1

mapped:
  key1: 100
  key2: 200

b:
  value: "subconfig"
  unexpected: "field"
  
env:
  - "PATH=/usr/bin"
  - "HOME=/home/user"

x: 10

unexpected1: true
unexpected2:
  v1: "123"
  v2: abc

`

	cfg := &exampleConfig{}
	cfg.rootConfig.unknownField = 999

	unknownFields, err := LoadConfigObject([]byte(exampleYamlContent), cfg)
	if err != nil {
		err = fmt.Errorf("加载配置失败: %+w", err)
		t.Fatalf("%+v", err)
	}

	spew.Dump(cfg)

	assert.Equal(t, "hello", cfg.stringVal)
	assert.Contains(t, cfg.stringsVal, "world")
	assert.Equal(t, 42, cfg.numberVal)
	assert.Contains(t, cfg.numbersVal, 43)
	assert.Contains(t, cfg.numbersVal, 44)
	assert.Equal(t, true, cfg.boolVal)
	assert.Equal(t, "subconfig", cfg.subConfig.stringVal)
	assert.Equal(t, 10, cfg.rootConfig.x)
	assert.Equal(t, 999, cfg.rootConfig.unknownField)
	assert.Equal(t, uint32(100), cfg.mapped["key1"])
	assert.Equal(t, uint32(200), cfg.mapped["key2"])
	assert.Equal(t, "/usr/bin", cfg.customLoader.values["PATH"])
	assert.Equal(t, "/home/user", cfg.customLoader.values["HOME"])
	assert.Contains(t, unknownFields, "unexpected1")
	assert.Contains(t, unknownFields, "unexpected2")
	assert.Contains(t, unknownFields, "b.unexpected")
}
