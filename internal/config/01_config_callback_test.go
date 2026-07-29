package config

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	log.SetOutput(t.Output())

	json := `{"test": "hello"}`
	var r test_main
	LoadConfigContent(json, &r)

	assert.Equal(t, "Receive:hello", r.Lvl2.Lvl3.Value)
}

func TestValidator(t *testing.T) {
	log.SetOutput(t.Output())

	json := `{"test": "hello"}`
	var r test_main
	LoadConfigContent(json, &r)

	// assert.Equal(t, "hello", r.Lvl1.Lvl2.Value)
	// assert.Equal(t, true, r.Validated)
	// assert.Equal(t, true, r.Lvl1.Validated)
	// assert.Equal(t, true, r.Lvl1.Lvl2.Validated)
}

type test_level1 struct {
	Value     string
	Validated bool
}

func (t *test_level1) FromString(v string) error {
	t.Value = "Receive:" + v
	return nil
}

func (t *test_level1) Validate() error {
	t.Validated = true
	return nil
}

type test_level2 struct {
	Lvl3      test_level1 `config:"test"`
	Validated bool
}

func (t *test_level2) Validate() error {
	t.Validated = true
	return nil
}

type test_main struct {
	Lvl2      *test_level2
	Validated bool
}

func (t *test_main) Validate() error {
	t.Validated = true
	return nil
}
