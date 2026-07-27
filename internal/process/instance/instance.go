package instance

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type ProcessInstance struct {
	internal *exec.Cmd
	err      error

	done chan struct{}

	everStarted bool

	beforeStart func(*exec.Cmd)
}

func New(cmdline []string) *ProcessInstance {
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Pdeathsig = syscall.SIGKILL

	return &ProcessInstance{
		internal: cmd,
		done:     make(chan struct{}),
	}
}

func (mc *ProcessInstance) SetSysProcAttr(setter func(attr *syscall.SysProcAttr)) {
	mc._assert_not_started()
	setter(mc.internal.SysProcAttr)
}

func (mc *ProcessInstance) SetEnv(envs []string) {
	mc._assert_not_started()
	mc.internal.Env = envs
}

func (mc *ProcessInstance) SetDir(dir string) {
	mc._assert_not_started()
	if dir == "" {
		var err error
		mc.internal.Dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		mc.internal.Dir = dir
	}
}

func (mc *ProcessInstance) SetBeforeStartHook(hook func(*exec.Cmd)) {
	mc._assert_not_started()
	if mc.beforeStart != nil {
		log.Println("警告: beforeStart hook已经被设置，新的hook将覆盖旧的hook")
	}
	mc.beforeStart = hook
}

func (mc *ProcessInstance) _assert_not_started() {
	if mc.everStarted {
		panic("命令已经启动，不能再修改启动信息")
	}
}

func (mc *ProcessInstance) IsStarted() bool {
	return mc.everStarted
}

func (mc *ProcessInstance) IsRunning() bool {
	if !mc.everStarted {
		return false
	}

	if mc.internal.ProcessState != nil || mc.err == nil { // 有这个说明wait返回了
		return false
	}

	return true
}

func (mc *ProcessInstance) Start() error {
	mc.everStarted = true // 这个程序里应该没有对这个有异步访问的可能性

	if mc.beforeStart != nil {
		mc.beforeStart(mc.internal)
	}

	// log.Printf("启动进程: %v", mc.internal.Args)
	err := mc.internal.Start()
	if err != nil {
		mc.err = err
		return err
	}

	go func() {
		mc.err = mc.internal.Wait()
		close(mc.done)
	}()

	return nil
}

func (mc *ProcessInstance) IsExited() bool {
	select {
	case <-mc.done:
		return true
	default:
		return false
	}
}

// Join 阻塞等待进程退出，并返回退出状态和错误信息
// 如果还没有启动，或者还在运行中，则会阻塞等待
func (mc *ProcessInstance) Join() (*os.ProcessState, error) {
	<-mc.done
	return mc.internal.ProcessState, mc.err
}

func (mc *ProcessInstance) Kill(signal os.Signal) error {
	return mc.internal.Process.Signal(signal)
}

func (mc *ProcessInstance) GetPid() int {
	if mc.internal.Process == nil {
		return 0
	}
	return mc.internal.Process.Pid
}
