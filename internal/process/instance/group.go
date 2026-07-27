package instance

import (
	_ "embed"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
)

//go:embed startup_script.sh
var startupScript string

type ProcessGroup struct {
	child_processes *ProcessList
	leaderInstance  *ProcessInstance

	env         []string
	cwd         string
	beforeStart func(*exec.Cmd)
	overlayRoot string

	mu sync.Mutex
}

func NewProcessGroup() *ProcessGroup {
	return &ProcessGroup{
		child_processes: NewProcessList(),
	}
}

func (pg *ProcessGroup) SetEnv(envs []string) {
	pg._assert_not_started()
	pg.env = envs
}

func (pg *ProcessGroup) SetDir(dir string) {
	pg._assert_not_started()
	pg.cwd = dir
}

func (pg *ProcessGroup) SetBeforeStartHook(hook func(*exec.Cmd)) {
	pg._assert_not_started()
	pg.beforeStart = hook
}

func (pg *ProcessGroup) SetOverlayRoot(root string) {
	pg._assert_not_started()
	pg.overlayRoot = root
}

func (pg *ProcessGroup) _assert_not_started() {
	if pg.leaderInstance != nil {
		panic("进程组已经启动，不能再修改启动信息")
	}
}

// 创建一个新的进程实例，并将其注册到进程组中，但不启动它
func (pg *ProcessGroup) CreateProcess(cmdline []string) *ProcessInstance {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	isLeader := pg.leaderInstance == nil

	var instance *ProcessInstance
	if isLeader {
		// TODO: 用golang实现这个helper
		instance = NewProcessInstance(append([]string{"bash", "-c", startupScript, "--"}, cmdline...))
		instance.SetEnv(append(pg.env, "_CHDIR_="+pg.cwd, "OVERLAY_ROOT="+pg.overlayRoot))
		instance.SetSysProcAttr(func(attr *syscall.SysProcAttr) {
			attr.Cloneflags |= syscall.CLONE_NEWNS | syscall.CLONE_NEWPID
		})

		pg.leaderInstance = instance
	} else {
		wrap := nsenter(pg.leaderInstance.GetPid(), pg.cwd, cmdline)
		// log.Printf("创建子进程: %v", wrap)
		instance = NewProcessInstance(wrap)
		instance.SetEnv(pg.env)
	}
	pg.child_processes.Register(instance)

	instance.SetDir("/")
	instance.SetBeforeStartHook(pg.beforeStart)

	return instance
}

func nsenter(leaderPid int, cwd string, cmdline []string) []string {
	return append([]string{"nsenter", "-t", strconv.Itoa(leaderPid), "--all", "--wdns=" + cwd, "--"}, cmdline...)
}
