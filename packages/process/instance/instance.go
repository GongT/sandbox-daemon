package instance

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/gongt/sandbox-daemon/packages/logger"
	"github.com/pkg/errors"
)

type ProcessInstance struct {
	internal *exec.Cmd
	err      error

	done        chan struct{}
	everStarted bool

	beforeStart func(*exec.Cmd)

	killer func(*ProcessInstance) error
}

func New(cmdline []string) *ProcessInstance {
	if len(cmdline) == 0 {
		panic("命令行不能为空")
	}
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

func defaultProcessKiller(process *ProcessInstance) error {
	// 此函数不可使用私有接口
	process.Kill(syscall.SIGTERM)

	select {
	case <-process.Wait():
		return nil
	case <-time.After(5 * time.Second):
		process.Kill(syscall.SIGKILL)
	}
	return nil
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

// 设置一个自定义的Stop方法，以实现优雅地关闭进程
func (mc *ProcessInstance) SetKiller(killer func(*ProcessInstance) error) {
	mc._assert_not_started()
	if mc.killer != nil {
		log.Println("警告: killer已经被设置，新的killer将覆盖旧的killer")
	}
	mc.killer = killer
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

	if mc.internal.ProcessState != nil || mc.err != nil { // 有这个说明wait返回了
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
		if err.Error() == "exec: already started" {
			return errors.WithMessage(err, "重复调用Start()")
		}

		mc.err = err
		close(mc.done)
		return errors.WithStack(err)
	}

	go func() {
		mc.err = mc.internal.Wait()
		logger.DProcF("ProcessInstance: 进程[%s]已退出", mc.String())
		close(mc.done)
	}()

	return nil
}

func (mc *ProcessInstance) GetResult() (*os.ProcessState, error) {
	if mc.internal.ProcessState == nil {
		return nil, errors.New("进程还没有退出（或从未启动），无法获取退出状态")
	}
	return mc.internal.ProcessState, nil
}

func (mc *ProcessInstance) MustResult() *os.ProcessState {
	if mc.internal.ProcessState == nil {
		panic("进程还没有退出（或从未启动），无法获取退出状态")
	}
	return mc.internal.ProcessState
}

func (mc *ProcessInstance) String() string {
	if mc.everStarted {
		if mc.IsExited() {
			return fmt.Sprintf("进程%d退出码%d", mc.GetPid(), mc.internal.ProcessState.ExitCode())
		} else {
			return fmt.Sprintf("进程%d", mc.GetPid())
		}
	} else {
		return fmt.Sprintf("未启动%v命令", mc.internal.Args[0])
	}
}

func (mc *ProcessInstance) Stop() error {
	defer mc.Join()
	if mc.killer != nil {
		return mc.killer(mc)
	} else {
		return defaultProcessKiller(mc)
	}
}

func (mc *ProcessInstance) IsExited() bool {
	select {
	case <-mc.done:
		return true
	default:
		return false
	}
}

// 阻塞等待进程退出，并返回退出状态和错误信息
// 如果还没有启动，或者还在运行中，则会阻塞等待直到进程（启动并）退出
func (mc *ProcessInstance) Join() (*os.ProcessState, error) {
	<-mc.done
	return mc.internal.ProcessState, mc.err
}

// 返回一个只读的 channel，实际不会读出任何值，当进程退出会关闭
func (mc *ProcessInstance) Wait() <-chan struct{} {
	return mc.done
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
