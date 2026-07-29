package reflection

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Field int
}

func TestIsTypeScalar(t *testing.T) {
	log.SetOutput(t.Output())

	require.True(t, IsTypeScalar(reflect.TypeFor[int]()))
	require.False(t, IsTypeScalar(reflect.TypeFor[[]int]()))
	require.False(t, IsTypeScalar(reflect.TypeFor[testStruct]()))
}

func TestIsTypeContainer(t *testing.T) {
	log.SetOutput(t.Output())

	require.True(t, IsTypeContainer(reflect.TypeFor[[]int]()))
	require.True(t, IsTypeContainer(reflect.TypeFor[map[int]string]()))
	require.False(t, IsTypeContainer(reflect.TypeFor[*int]()))
	require.False(t, IsTypeContainer(reflect.TypeFor[int]()))
	require.False(t, IsTypeContainer(reflect.TypeFor[testStruct]()))
}

func TestIsTypeSerializable(t *testing.T) {
	log.SetOutput(t.Output())

	require.True(t, IsTypeSerializable(reflect.TypeFor[int]()))
	require.False(t, IsTypeSerializable(reflect.TypeFor[*int]()))
}

func TestIndirect(t *testing.T) {
	log.SetOutput(t.Output())

	var a ****int
	require.Equal(t, reflect.TypeFor[int](), IndirectType(reflect.TypeOf(a)))

	v := IndirectValue(reflect.ValueOf(a))
	require.False(t, v.IsValid())

	var b int = 42
	var c *int = &b
	var d **int = &c

	v2 := IndirectValue(reflect.ValueOf(d))
	require.Equal(t, 42, v2.Interface())
}

func TestInstantiate(t *testing.T) {
	log.SetOutput(t.Output())

	var a ****int
	val1 := InstantiateType(reflect.TypeOf(a))
	require.Equal(t, 0, val1.Interface())

	type sub struct {
		field int
	}
	type deepStruct struct {
		mapField    map[string]int
		simpleField string
	}

	var s *deepStruct
	val2ptr := InstantiateType(reflect.TypeOf(s))
	val2 := val2ptr.Interface().(deepStruct)

	require.Nil(t, s)
	require.NotNil(t, val2)
	require.Equal(t, "", val2.simpleField)
	require.Nil(t, val2.mapField)
}

func TestPointers(t *testing.T) {
	log.SetOutput(t.Output())

	var a ****string

	final, err := InstantiatePointers(reflect.ValueOf(&a))

	require.NoError(t, err)
	require.Nil(t, final.Interface())
	str := "hello"
	final.Set(reflect.ValueOf(&str))
	require.Equal(t, "hello", ****a)
}
