package convert

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type canString struct {
	extra string
	value string
}

func (c *canString) String() string {
	return "I'm " + c.value
}

func (c *canString) FromString(s string) {
	c.value = s
}

func TestConvertString(t *testing.T) {
	log.SetOutput(t.Output())

	var res1 string
	res1, err := ConvertToString(123)
	require.NoError(t, err)
	assert.Equal(t, "123", res1)

	res1, err = ConvertToString("hello")
	require.NoError(t, err)
	assert.Equal(t, "hello", res1)

	res1, err = ConvertToString(map[string]int{"a": 1})
	require.Error(t, err)

	res1, err = ConvertToString(canString{value: "test"})
	require.NoError(t, err)
	assert.Equal(t, "I'm test", res1)

	var res2 int
	err = ConvertStringTo("456", &res2)
	require.NoError(t, err)
	assert.Equal(t, 456, res2)

	var res3 float64
	err = ConvertStringTo("3.14", &res3)
	require.NoError(t, err)
	assert.Equal(t, 3.14, res3)

	res4 := canString{
		value: "unknown", extra: "unchanged",
	}
	err = ConvertStringTo("world", &res4)
	require.NoError(t, err)
	assert.Equal(t, "I'm world", res4.String())
	assert.Equal(t, "unchanged", res4.extra)
}
