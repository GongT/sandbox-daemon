package process_list

import (
	"errors"
	"slices"
	"sync"
)

// ProcessList 是一个进程列表，里面的元素是 *ProcessInstance
// 用于管理列表中进程的退出顺序（即先启动的进程后退出）
type ProcessList struct {
	instances []stoppable
	mu        sync.RWMutex
}

type stoppable interface {
	Wait() <-chan struct{}
	Stop() error
}

func New() *ProcessList {
	r := &ProcessList{
		instances: make([]stoppable, 0),
		mu:        sync.RWMutex{},
	}
	return r
}

func (pl *ProcessList) Register(instance stoppable) int {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.instances = append(pl.instances, instance)

	go func() {
		<-instance.Wait()
		pl.unregister(instance)
	}()

	return len(pl.instances)
}

func (pl *ProcessList) unregister(instance stoppable) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for i, inst := range pl.instances {
		if inst == instance {
			pl.instances = append(pl.instances[:i], pl.instances[i+1:]...)
			break
		}
	}
}

func (pl *ProcessList) StopAll() error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	var errs []error
	for _, instance := range slices.Backward(pl.instances) {
		err := instance.Stop()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (pl *ProcessList) ReadAccess(access func([]stoppable)) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	access(pl.instances)
}

func (pl *ProcessList) Size() int {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	return len(pl.instances)
}
