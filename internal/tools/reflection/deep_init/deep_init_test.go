package deep_init

import (
	"log"
	"testing"

	"github.com/goforj/godump"
	"github.com/gongt/sandbox-daemon/internal/tools/interfaces"
	"github.com/stretchr/testify/require"
)

type fill_sub struct {
	Field1          string
	Field2          int
	ChannelField    chan int
	privateField    bool
	unexportedField []string
}

func (v *fill_sub) Initialize() {
	v.privateField = true
}

type fill_main struct {
	SubField    fill_sub
	SubFieldPtr *fill_sub
	MapField    map[string]string
	ScalarField string
	SliceField  []int
}

func TestFill(t *testing.T) {
	log.SetOutput(t.Output())

	var test fill_main
	godump.Fdump(t.Output(), test)

	walkedptrs := DeepInitialize(&test)

	for _, ptr := range walkedptrs {
		if sub, ok := ptr.(interfaces.Initializer); ok {
			sub.Initialize()
		}
	}

	godump.Fdump(t.Output(), test)

	require.Equal(t, true, test.SubField.privateField)
	require.Nil(t, test.SubField.ChannelField)
	require.Nil(t, test.SubField.unexportedField)

	require.NotNil(t, test.SubFieldPtr)
	require.Equal(t, true, test.SubFieldPtr.privateField)

	require.NotNil(t, test.MapField)
	require.NotNil(t, test.SliceField)
}
