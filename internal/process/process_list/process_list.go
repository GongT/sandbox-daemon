package process_list

import (
	"fmt"
	"sync"

	"github.com/gongt/sandbox-daemon/internal/tools"
	"gitlab.com/tozd/go/errors"
)

// ProcessList 是一个进程列表，里面的元素是 *ProcessInstance
// 用于管理列表中进程的退出顺序（即先启动的进程后退出）
type ProcessList struct {
	instances []*processBox

	mu sync.RWMutex
}

type processBox struct {
	instance stoppable
	notify   chan struct{}
}

type stoppable interface {
	Wait() <-chan struct{}
	Stop() error
}

func New() *ProcessList {
	r := &ProcessList{
		instances: make([]*processBox, 0),
		mu:        sync.RWMutex{},
	}
	return r
}

func (pl *ProcessList) Register(instance stoppable) int {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	tools.DebugLog("ProcessList.Register: 注册进程 %s (当前%d个)", maybe_string(instance), len(pl.instances))
	box := &processBox{
		instance: instance,
		notify:   make(chan struct{}),
	}
	pl.instances = append(pl.instances, box)

	go func() {
		<-instance.Wait()
		pl.unregister(instance)
		close(box.notify)
	}()

	return len(pl.instances)
}

func (pl *ProcessList) unregister(instance stoppable) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for i, inst := range pl.instances {
		if inst.instance == instance {
			pl.instances = append(pl.instances[:i], pl.instances[i+1:]...)
			break
		}
	}
}

func (pl *ProcessList) StopAll() error {
	var errs []error
	peek := func() *processBox {
		pl.mu.RLock()
		defer pl.mu.RUnlock()

		if len(pl.instances) == 0 {
			return nil
		}
		return pl.instances[len(pl.instances)-1]

	}
	for {
		box := peek()
		if box == nil {
			tools.DebugLog("ProcessList.StopAll: 已全部停止")
			break
		}

		tools.DebugLog("ProcessList.StopAll: 停止进程 %s", maybe_string(box.instance))

		err := box.instance.Stop()

		if err != nil {
			errs = append(errs, err)
		}

		<-box.notify

		tools.DebugLog("ProcessList.StopAll: 进程 %s 已停止", maybe_string(box.instance))
	}

	return errors.Join(errs...)
}

func (pl *ProcessList) Size() int {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	return len(pl.instances)
}

func maybe_string(instance stoppable) string {
	if s, ok := instance.(fmt.Stringer); ok {
		return s.String()
	}
	return "<?>"
}
