package process_list

import (
	"sync"

	"github.com/gongt/sandbox-daemon/internal/process/instance"
)

// ProcessList 是一个进程列表，里面的元素是 *ProcessInstance
type ProcessList struct {
	instances []*instance.ProcessInstance
	mu        sync.RWMutex
}

func New() *ProcessList {
	return &ProcessList{
		instances: make([]*instance.ProcessInstance, 0),
		mu:        sync.RWMutex{},
	}
}

func (pl *ProcessList) Register(instance *instance.ProcessInstance) int {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.instances = append(pl.instances, instance)

	go func() {
		// 等待进程退出
		instance.Join()

		// 进程退出后，从列表中移除
		pl.mu.Lock()
		defer pl.mu.Unlock()

		for i, inst := range pl.instances {
			if inst == instance {
				pl.instances = append(pl.instances[:i], pl.instances[i+1:]...)
				break
			}
		}
	}()

	return len(pl.instances)
}

func (pl *ProcessList) ReadAccess(access func([]*instance.ProcessInstance)) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	access(pl.instances)
}
