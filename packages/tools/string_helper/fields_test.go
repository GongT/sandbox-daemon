package string_helper

import (
	"iter"
	"testing"

	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCutFields(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	s := "a,b;c,d"
	seps := []string{".", ",", ";"}
	itr := CutFields(s, seps)

	next, stop := iter.Pull2(itr)
	defer stop()

	var ok bool
	var value string
	var sep int

	value, sep, ok = next()
	require.True(t, ok)
	assert.Equal(t, "a", value)
	assert.Equal(t, 1, sep)

	value, sep, ok = next()
	require.True(t, ok)
	assert.Equal(t, "b", value)
	assert.Equal(t, 2, sep)

	value, sep, ok = next()
	require.True(t, ok)
	assert.Equal(t, "c", value)
	assert.Equal(t, 1, sep)

	value, sep, ok = next()
	require.True(t, ok)
	assert.Equal(t, "d", value)
	assert.Equal(t, -1, sep)

	value, sep, ok = next()
	require.False(t, ok)
}
