package reflection

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAssign(t *testing.T) {
	log.SetOutput(t.Output())

	var a int

	err := AssignValue(&a, 42)
	require.NoError(t, err)
	require.Equal(t, 42, a)

	m := map[string]int{
		"OLD_VAL": 1,
	}

	err = AssignValue(&m, map[string]int{"a": 1})
	require.NoError(t, err)
	require.Equal(t, map[string]int{"a": 1}, m)
}
