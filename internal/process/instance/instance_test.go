package instance

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInstance(t *testing.T) {
	proc := NewProcessInstance([]string{"echo", "hello"})
	require.Equal(t, 0, proc.GetPid())

	require.NoError(t, proc.Start())
	require.Error(t, proc.Start())

	stat, err := proc.Join()
	require.NoError(t, err)
	require.Equal(t, 0, stat.ExitCode())

	stat, err = proc.Join()
	require.NoError(t, err)
	require.Equal(t, 0, stat.ExitCode())
}

func TestInstanceStop(t *testing.T) {
	proc := NewProcessInstance([]string{"sleep", "infinity"})
	require.NoError(t, proc.Start())

	go func() {
		// 等待一段时间后停止进程
		time.Sleep(time.Second / 2)
		proc.Kill(os.Kill)
	}()

	stat, err := proc.Join()
	require.Error(t, err)
	require.Equal(t, -1, stat.ExitCode())

	require.Error(t, proc.Kill(os.Kill))

	fmt.Printf("退出状态: %+v\n", stat.String())
}
