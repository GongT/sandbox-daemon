package daemon

import (
	"log"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/instance"
	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/pkg/errors"
)

var d *instance.D

type InitCommand struct {
	internal.WithSessionCommand
}

func (config *InitCommand) Run(runtime *myenv.GlobalOptions) error {
	if d != nil {
		panic("DaemonInstance: 只能启动一个守护进程实例")
	}

	d = instance.New(&config.WithSessionCommand, runtime)

	log.Printf("守护进程已启动，session_id: %s", config.SessionId.String())
	d.Join()

	err := d.Destroy()
	if err != nil {
		return errors.Wrap(err, "退出时出现问题")
	}

	return nil
}

func GetDaemon() *instance.D {
	return d
}

func Destroy() error {
	err := d.Destroy()
	d = nil
	return err
}
