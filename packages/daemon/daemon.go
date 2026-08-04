package daemon

import (
	"errors"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/reap"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/rpc"
)

type DaemonInstance struct {
	mu sync.RWMutex

	parts []internal.DaemonComponent
}

var instance *DaemonInstance

func StartDaemon() *DaemonInstance {
	if instance != nil {
		panic("DaemonInstance: 只能启动一个守护进程实例")
	}
	instance = &DaemonInstance{}

	// 防止stdout、stderr被关闭时，程序直接退出
	signal.Reset(syscall.SIGPIPE)

	rpcServer := rpc.NewServer()
	reaper := reap.New()
	mp := main_process.New()

	instance.parts = []internal.DaemonComponent{
		rpcServer,
		reaper,
		mp,
	}

	return instance
}

func GetDaemon() *DaemonInstance {
	return instance
}

func (pm *DaemonInstance) Destroy() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	errs := []error{}
	for _, part := range pm.parts {
		err := part.Stop()
		if err != nil {
			errs = append(errs, err)
		}
	}

	instance = nil

	return errors.Join(errs...)
}
