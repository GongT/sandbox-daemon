package environ

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvMap(t *testing.T) {
	var err error
	env := Map{}

	env.SetLine("FOO=BAR")
	require.Equal(t, "BAR", env["FOO"])

	env.SetLine("啊=QUX")
	require.Equal(t, "QUX", env["啊"])

	err = env.SetLine("INVALID_LINE")
	require.Error(t, err)

	env.Clear()
	require.Equal(t, 0, env.Size())

	env.Extend(map[string]string{
		"A": "1",
		"B": "2",
	}, true)
	require.Equal(t, 2, env.Size())
	require.Equal(t, "1", env["A"])
	require.Equal(t, "2", env["B"])

	env.Extend(map[string]string{"A": "3"}, true)
	require.Equal(t, "3", env["A"])

	env.Extend(map[string]string{"A": "4"}, false)
	require.Equal(t, "3", env["A"])

	env.Delete("B")
	require.Equal(t, 1, env.Size())
	require.False(t, env.Has("B"))

	env.Delete("C")

	env.ExtendLines([]string{"X=Y", "Z=W", "Q"}, true)
	require.Equal(t, 3, env.Size())

	lines := env.ToLines()
	require.Equal(t, 3, len(lines))
	require.Contains(t, lines, "A=3")
	require.Contains(t, lines, "X=Y")
	require.Contains(t, lines, "Z=W")
}
