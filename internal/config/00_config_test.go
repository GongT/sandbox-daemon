package config

import (
	"log"
	"testing"

	"github.com/goforj/godump"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfigPart1 struct {
	Name        string            `config:"name"`
	Enabled     bool              `config:"enabled"`
	Commandline []string          `config:"commandline"`
	Environment map[string]string `config:"environment"`
	// Timeout     time.Duration     `config:"timeoutMs"`
	// Urls        []url.URL         `config:"urls"`

	BuildConfig testDeep1

	UnrelatedField string
}

type testConfigPart2 struct {
	Objects []testObject      `config:"objects"`
	MapA    mapAlias          `config:"map_field"`
	MapB    map[string]string `config:"map_field"`
}

type testDeep1 struct {
	Deep testDeep2
}

type testDeep2 struct {
	Value string `config:"not.really.deep"`
}

const testConfigContentYaml = `
name: "value-of-name"
enabled: true
commandline:
  - "arg1"
  - "arg2"
environment:
  Key1: "value1"
  key2: "value2"
timeoutMs: 1000
urls: "https://example.com/path?query=1"
hidden_field: "should not be loaded"
objects:
  - id: 1
    name: "object1"
  - id: 2
    name: "object2"
map_field:
  keyA: "valueA"
  keyB: "valueB"

not:
  really:
    deep: "load-to-deep"
`

const testConfigContentJson = `
name: "value-of-name"
enabled: true
commandline:
  - "arg1"
  - "arg2"
environment:
  Key1: "value1"
  key2: "value2"
timeoutMs: 1000
urls: "https://example.com/path?query=1"
hidden_field: "should not be loaded"
objects:
  - id: 1
    name: "object1"
  - id: 2
    name: "object2"
map_field:
  keyA: "valueA"
  keyB: "valueB"

not:
  really:
    deep: "load-to-deep"
`

func TestLoadYaml(t *testing.T) {
	log.SetOutput(t.Output())

	testData(t, testConfigContentYaml)
}

func TestLoadJson(t *testing.T) {
	log.SetOutput(t.Output())

	testData(t, testConfigContentJson)
}

func testData(t *testing.T, content string) {
	var part2 testConfigPart2
	var part1 testConfigPart1

	part1 = testConfigPart1{
		Environment: map[string]string{
			"HELLO": "WORLD",
		},
	}

	err := LoadConfigContent(content, &part1, &part2)
	require.NoError(t, err)

	godump.Fdump(t.Output(), part1, part2)

	assert.Equal(t, "value-of-name", part1.Name)
	assert.Equal(t, true, part1.Enabled)
	assert.Equal(t, []string{"arg1", "arg2"}, part1.Commandline)
	assert.Equal(t, map[string]string{"HELLO": "WORLD", "Key1": "value1", "key2": "value2"}, part1.Environment)
	// assert.Equal(t, time.Duration(1000)*time.Millisecond, part1.Timeout)
	// require.Len(t, part1.Urls, 1)
	// assert.Equal(t, "https://example.com/path?query=1", part1.Urls[0].String())

	require.Len(t, part2.Objects, 2)
	assert.Equal(t, 1, part2.Objects[0].ID)
	assert.Equal(t, "object1", part2.Objects[0].Name)
	assert.Equal(t, 2, part2.Objects[1].ID)
	assert.Equal(t, "object2", part2.Objects[1].Name)

	assert.Equal(t, mapAlias{"keyA": "valueA", "keyB": "valueB"}, part2.MapA)
	assert.Equal(t, map[string]string{"keyA": "valueA", "keyB": "valueB"}, part2.MapB)

	assert.Equal(t, "load-to-deep", part1.BuildConfig.Deep.Value)
}

type testNoHiddenFieldHasTag struct {
	hiddenField string `config:"hidden_field"`
}

func TestNoHiddenFieldHasTag(t *testing.T) {
	log.SetOutput(t.Output())

	var part testNoHiddenFieldHasTag
	err := LoadConfigContent("", &part)
	require.Error(t, err, "字段hiddenField是未导出的，但有config标签")
}
