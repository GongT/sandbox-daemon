package daemon

import (
	"log"
	"os"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	daemon_config "github.com/gongt/sandbox-daemon/packages/daemon/internal/config"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process/mp_config"
	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/pkg/errors"
)

type StartCommand struct {
	internal.WithSessionCommand

	configFile daemon_config.ConfigFileOptionOptional `long:"config" description:"设置配置文件路径"`
	config     *daemon_config.DaemonConfig
}

func (cmd *StartCommand) Run(runtime *myenv.GlobalOptions) error {
	if err := GetDaemon().Assert(); err != nil {
		return err
	}

	if cmd.configFile.IsEmpty() {
		return errors.New("请指定配置文件路径")
	}

	bs, err := os.ReadFile(cmd.configFile.GetFilePath())
	if err != nil {
		return errors.WithMessagef(err, "无法读取配置文件: %s", cmd.configFile.GetFilePath())
	}

	text := string(bs)

	mp, err := mp_config.New(text)
	if err != nil {
		return errors.WithMessage(err, "解读配置主进程信息失败")
	}
	if !mp.HasExecSpecified() {
		return errors.New("配置文件没有指定主进程信息")
	}

	log.Println("配置文件存在主进程信息，立即启动……")
	return GetDaemon().LaunchMainProcess(mp)
}
