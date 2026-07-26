package daemon

import (
	"log"
	"sync"

	reap "github.com/hashicorp/go-reap"
)

type DaemonInstance struct {
	mu   *sync.RWMutex
	done chan struct{}
}

var single_instance bool

func StartDaemon() *DaemonInstance {
	if single_instance {
		panic("DaemonInstance: 只能启动一个守护进程实例")
	}
	single_instance = true

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
	}
}

func (pm *DaemonInstance) Stop() {
	close(pm.done)
	single_instance = false
}
