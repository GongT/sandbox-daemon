package daemon

import (
	"log"
	"os"

	"github.com/gongt/sandbox-daemon/packages/config_loader"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	daemon_config "github.com/gongt/sandbox-daemon/packages/daemon/internal/config"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/instance"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process/mp_config"
	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/gongt/sandbox-daemon/packages/tools/signals"
	"github.com/pkg/errors"
)

var d *instance.D

type InitCommand struct {
	internal.WithSessionCommand

	ConfigFileSpec daemon_config.ConfigFileOption `long:"config" description:"设置配置文件路径" default-mask:"/etc/sandbox-daemon.yaml"`
	initConfigText string
	config         *daemon_config.DaemonConfig
}

func (cmd *InitCommand) Run(runtime *myenv.GlobalOptions) error {
	if d.Assert() == nil {
		return errors.New("守护进程实例已存在")
	}

	if !cmd.ConfigFileSpec.IsEmpty() {
		bs, _ := os.ReadFile(cmd.ConfigFileSpec.GetFilePath())
		if bs != nil {
			cmd.initConfigText = string(bs)
		}
		if cmd.initConfigText != "" {
			cmd.config = &daemon_config.DaemonConfig{}
			err := config_loader.LoadConfigContent(cmd.initConfigText, cmd.config)
			if err != nil {
				return errors.WithMessage(err, "解读配置文件失败")
			}
		}
	}

	d = instance.New(&cmd.WithSessionCommand, cmd.config)

	log.Printf("守护进程已启动，session_id: %s", cmd.SessionId.String())

	if cmd.initConfigText != "" {
		mp, err := mp_config.New(cmd.initConfigText)
		if err != nil {
			return errors.WithMessage(err, "解读配置主进程信息失败")
		}
		if mp.HasExecSpecified() {
			log.Println("配置文件存在主进程信息，立即启动……")
			err = d.LaunchMainProcess(mp)
			if err != nil {
				return errors.WithMessage(err, "启动主进程失败")
			}
		} else {
			log.Println("配置文件没有指定主进程信息.")
		}
	} else {
		log.Println("配置文件不存在.")
	}

	<-d.Join()
	log.Println("守护进程退出，session_id: ", cmd.SessionId.String())

	err := d.Destroy()

	signals.AppQuit.Set(0)

	if err != nil {
		log.Println(errors.WithMessage(err, "(忽略)守护进程退出时出现问题"))
	}

	d = nil
	return nil
}

func GetDaemon() *instance.D {
	return d
}
