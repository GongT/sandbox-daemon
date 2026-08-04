package process_group

import (
	_ "embed"
	"errors"
	"os/exec"
	"strconv"
	"sync"
	"syscall"

	"github.com/gongt/sandbox-daemon/packages/process/instance"
	"github.com/gongt/sandbox-daemon/packages/process/process_list"
)

//go:embed startup_script.sh
var startupScript string

type ProcessGroup struct {
	child_processes *process_list.ProcessList
	leaderInstance  *instance.ProcessInstance

	env         []string
	cwd         string
	beforeStart func(*exec.Cmd)
	config      NamespaceConfig

	mu sync.RWMutex
}

func New() *ProcessGroup {
	return &ProcessGroup{
		child_processes: process_list.New(),
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
	pg.config.OverlayRoot = root
}

func (pg *ProcessGroup) _assert_not_started() {
	pg.mu.RLock()
	defer pg.mu.RUnlock()

	if pg.leaderInstance != nil {
		panic("进程组已经启动，不能再修改启动信息")
	}
}

func (pg *ProcessGroup) StopAll() error {
	return pg.child_processes.StopAll()
}

// 创建一个新的进程实例，并将其注册到进程组中，但不启动它
func (pg *ProcessGroup) CreateProcess(cmdline []string) (*instance.ProcessInstance, error) {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	isLeader := pg.leaderInstance == nil

	var proc *instance.ProcessInstance
	if isLeader {
		// TODO: 用golang实现这个helper
		proc = instance.New(append([]string{"bash", "-c", startupScript, "--"}, cmdline...))
		proc.SetEnv(append(pg.env, "_CHDIR_="+pg.cwd, "OVERLAY_ROOT="+pg.config.OverlayRoot))
		proc.SetSysProcAttr(func(attr *syscall.SysProcAttr) {
			attr.Cloneflags |= syscall.CLONE_NEWNS | syscall.CLONE_NEWPID
		})

		pg.leaderInstance = proc
	} else {
		pid := pg.leaderInstance.GetPid()
		if pid == 0 {
			return nil, errors.New("leader进程还没有启动，无法创建子进程")
		}

		wrap := nsenter(pid, pg.cwd, cmdline)
		// log.Printf("创建子进程: %v", wrap)
		proc = instance.New(wrap)
		proc.SetEnv(pg.env)
	}
	pg.child_processes.Register(proc)

	proc.SetDir("/") // 使用了mountns，loader的工作目录没有意义
	proc.SetBeforeStartHook(pg.beforeStart)

	return proc, nil
}

func nsenter(leaderPid int, cwd string, cmdline []string) []string {
	return append([]string{"nsenter", "-t", strconv.Itoa(leaderPid), "--all", "--wdns=" + cwd, "--"}, cmdline...)
}
