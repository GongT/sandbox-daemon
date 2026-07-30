package instance

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInstance(t *testing.T) {
	log.SetOutput(t.Output())

	proc := New([]string{"echo", "hello"})
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
	log.SetOutput(t.Output())

	proc := New([]string{"sleep", "infinity"})
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

	log.Printf("退出状态: %+v", stat.String())
}

func (mc *ProcessInstance) waitForTest(t *testing.T) {
	select {
	case <-mc.done:
	case <-time.After(time.Second * 2):
		t.Fatal("等待进程退出超时")
	}
}

func TestErrorInstance(t *testing.T) {
	log.SetOutput(t.Output())

	proc := New([]string{"nonexistent_command"})
	require.Error(t, proc.Start())

	proc.waitForTest(t)
	stat, err := proc.Join()
	require.Error(t, err)
	require.Equal(t, -1, stat.ExitCode())
}
