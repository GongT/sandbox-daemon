package reap

import (
	"log"
	"sync"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	"github.com/hashicorp/go-reap"
)

var _ internal.DaemonComponent = (*processReaper)(nil)

type processReaper struct {
	mu   sync.RWMutex
	done chan struct{}

	pids   reap.PidCh
	errors reap.ErrorCh
}

func New() *processReaper {
	return &processReaper{
		mu:     sync.RWMutex{},
		done:   make(chan struct{}),
		pids:   make(reap.PidCh, 1),
		errors: make(reap.ErrorCh, 1),
	}
}

func (r *processReaper) Start() error { // 进程回收
	go func() {
		reap.ReapChildren(r.pids, r.errors, r.done, &r.mu)

		close(r.pids)
		close(r.errors)
	}()

	// debug
	go func() {
	loop:
		for {
			select {
			case pid := <-r.pids:
				log.Printf("进程退出: %d", pid)
			case err := <-r.errors:
				log.Printf("进程回收错误: %v", err)
			case <-r.done:
				break loop
			}
		}
	}()

	return nil
}

func (r *processReaper) Stop() error {
	close(r.done)
	return nil
}
