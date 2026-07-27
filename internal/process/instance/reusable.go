package instance

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type ReusableProcessInstance struct {
	instance *ProcessInstance

	env           []string
	cmdline       []string
	dir           string
	beforeStart   func(*exec.Cmd)
	sysProcSetter func(attr *syscall.SysProcAttr)
}

func NewReusable() *ReusableProcessInstance {
	return &ReusableProcessInstance{}
}

func (rpi *ReusableProcessInstance) SetCmdline(cmdline []string) {
	rpi.cmdline = cmdline
}

func (rpi *ReusableProcessInstance) SetEnv(envs []string) {
	rpi.env = envs
}

func (rpi *ReusableProcessInstance) SetSysProcAttr(setter func(attr *syscall.SysProcAttr)) {
	rpi.sysProcSetter = (setter)
}

func (rpi *ReusableProcessInstance) SetDir(dir string) {
	rpi.dir = dir
}

func (rpi *ReusableProcessInstance) SetBeforeStartHook(hook func(*exec.Cmd)) {
	rpi.beforeStart = hook
}

func (rpi *ReusableProcessInstance) IsStarted() bool {
	if rpi.instance == nil {
		return false
	}

	return rpi.instance.IsStarted()
}

func (rpi *ReusableProcessInstance) IsRunning() bool {
	if rpi.instance == nil {
		return false
	}

	return rpi.instance.IsRunning()
}

func (rpi *ReusableProcessInstance) create() {
	rpi.instance = New(rpi.cmdline)
	rpi.instance.SetEnv(rpi.env)
	rpi.instance.SetDir(rpi.dir)
	if rpi.beforeStart != nil {
		rpi.instance.SetBeforeStartHook(rpi.beforeStart)
	}
	if rpi.sysProcSetter != nil {
		rpi.instance.SetSysProcAttr(rpi.sysProcSetter)
	}
}

func (rpi *ReusableProcessInstance) Start() error {
	if rpi.instance == nil {
		rpi.create()
	}
	if rpi.instance.IsExited() {
		rpi.create()
	}

	return rpi.instance.Start()
}

func (rpi *ReusableProcessInstance) IsExited() bool {
	if rpi.instance == nil {
		return false
	}

	return rpi.instance.IsExited()
}

func (rpi *ReusableProcessInstance) Join() (*os.ProcessState, error) {
	if rpi.instance == nil {
		return nil, fmt.Errorf("reusableProcessInstance: 试图Join一个还没有启动的实例")
	}

	return rpi.instance.Join()
}

func (rpi *ReusableProcessInstance) Kill(signal os.Signal) error {
	return rpi.instance.Kill(signal)
}

func (rpi *ReusableProcessInstance) GetPid() int {
	if rpi.instance == nil {
		return 0
	}
	return rpi.instance.GetPid()
}
