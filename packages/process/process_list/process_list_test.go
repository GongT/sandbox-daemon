package process_list

import (
	"testing"
	"time"

	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/stretchr/testify/require"
)

type mockProcess struct {
	sleep time.Duration
	error error
	wait  chan struct{}
}

func mock(sleep time.Duration, error error) *mockProcess {
	r := &mockProcess{
		sleep: sleep,
		error: error,
		wait:  make(chan struct{}),
	}
	r.Start()
	return r
}

func (mp *mockProcess) String() string {
	if mp.error == nil {
		return "sleep(" + mp.sleep.String() + ")"
	} else {
		return "sleep(" + mp.sleep.String() + ", error)"
	}
}

func (mp *mockProcess) Start() {
	if mp.sleep == 0 {
		close(mp.wait)
		return
	}
	go func() {
		time.Sleep(mp.sleep)
		close(mp.wait)
	}()
}

func (mp *mockProcess) Wait() <-chan struct{} {
	return mp.wait
}

func (mp *mockProcess) Stop() error {
	return mp.error
}

func TestProcessList_Start(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	list1 := New()
	list1.Register(mock(1*time.Second, nil))
	require.Equal(t, 1, len(list1.instances))

	list1.Register(mock(time.Second/2, nil))
	require.Equal(t, 2, len(list1.instances))

	list1.Register(mock(0, nil))
	require.Equal(t, 3, len(list1.instances))

	time.Sleep(1 * time.Millisecond)
	require.Equal(t, 2, len(list1.instances))

	time.Sleep(1*time.Second + 200*time.Millisecond)
	require.Equal(t, 0, len(list1.instances))
}
