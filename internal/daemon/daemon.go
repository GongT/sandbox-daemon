package daemon

import (
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gongt/sandbox-daemon/internal/process/process_group"
	reap "github.com/hashicorp/go-reap"
)

type DaemonInstance struct {
	mu   *sync.RWMutex
	done chan struct{}

	pg *process_group.ProcessGroup
}

var single_instance bool

func StartDaemon() *DaemonInstance {
	if single_instance {
		panic("DaemonInstance: 只能启动一个守护进程实例")
	}
	single_instance = true

	// 防止stdout、stderr被关闭时，程序直接退出
	signal.Reset(syscall.SIGPIPE)

	// 进程回收
	pids := make(reap.PidCh, 1)
	errors := make(reap.ErrorCh, 1)
	done := make(chan struct{})
	mu := &sync.RWMutex{}

	go reap.ReapChildren(pids, errors, done, mu)

	go func() {
		for {
			select {
			case pid := <-pids:
				log.Printf("进程退出: %d", pid)
			case err := <-errors:
				log.Printf("进程回收错误: %v", err)
			case <-done:
				return
			}
		}
	}()

	return &DaemonInstance{
		mu:   mu,
		done: done,
		pg:   process_group.New(),
	}
}

func (pm *DaemonInstance) Destroy() {
	close(pm.done)
	single_instance = false
}

func (pm *DaemonInstance) StopMainProcess() error {
	return pm.pg.StopAll()
}
