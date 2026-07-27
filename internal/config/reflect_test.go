package config

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type reflectTestConfig struct {
	Name    string             `config:"name"`
	Enabled bool               `config:"enabled"`
	StrArr  []string           `config:"strings"`
	Objects []testObject       `config:"objects"`
	MapA    mapAlias           `config:"map"`
	MapB    map[string]string  `config:"map"`
	Nested  *reflectTestNested `config:"nested"`
	Age     *int               `config:"age"`
	Skipped chan int
}

type reflectTestNested struct {
	Value int `config:"value"`
}

type testObject struct {
	ID   int    `config:"id"`
	Name string `config:"name"`
}

type mapAlias map[string]string

type testCtx struct {
	ConfigFillContext

	t *testing.T
}

func (ctx *testCtx) HasValue(tagPath ConfigPath) (bool, error) {
	return true, nil
}

func (ctx *testCtx) GetArraySize(tagPath ConfigPath) (int, error) {
	require.Equal(ctx.t, tagPath.Size, 1)
	require.Contains(ctx.t, []string{"strings", "objects"}, tagPath.StringAt(0))
	return 2, nil
}

func (ctx *testCtx) GetObjectKeys(tagPath ConfigPath) ([]string, error) {
	require.Equal(ctx.t, tagPath.Size, 1)
	require.Equal(ctx.t, "map", tagPath.StringAt(0))
	return []string{"keyA", "keyB"}, nil
}

func (ctx *testCtx) GetValue(typ reflect.Type, tagPath ConfigPath) (string, error) {
	require.NotEmpty(ctx.t, tagPath)
	last := tagPath.StringAt(-1)
	switch last {
	case "keyA":
		return "valueA", nil
	case "keyB":
		return "valueB", nil
	}
	switch typ.Kind() {
	case reflect.String:
		return "value:" + last, nil
	case reflect.Bool:
		return "true", nil
	case reflect.Int:
		return "42", nil
	default:
		return "", nil
	}
}

func (ctx *testCtx) ConvertNonScalar(get_value string, t reflect.Type) (interface{}, error) {
	return nil, fmt.Errorf("没有此路径")
}

func TestWalkStruct(t *testing.T) {
	log.SetOutput(t.Output())

	input := &reflectTestConfig{}
	ctx := &testCtx{t: t}

	err := WalkStruct(input, ctx)
	require.NoError(t, err)

	require.Equal(t, "value:name", input.Name)
	require.True(t, input.Enabled)
	require.Len(t, input.StrArr, 2)
	for i, item := range input.StrArr {
		require.Equal(t, "value:"+strconv.Itoa(i), item)
	}

	require.Len(t, input.Objects, 2)
	for _, object := range input.Objects {
		require.Equal(t, 42, object.ID)
		require.Equal(t, "value:name", object.Name)
	}

	require.Equal(t, mapAlias{"keyA": "valueA", "keyB": "valueB"}, input.MapA)
	require.Equal(t, map[string]string{"keyA": "valueA", "keyB": "valueB"}, input.MapB)
	require.NotNil(t, input.Nested)
	require.Equal(t, 42, input.Nested.Value)
	require.NotNil(t, input.Age)
	require.Equal(t, 42, *input.Age)
	require.Nil(t, input.Skipped)
}
