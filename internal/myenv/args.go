package myenv

import (
	"os"
	"strings"

)

var RuntimeConfig *GlobalOptions

type GlobalOptions struct {
	RuntimeDir    string `long:"dir" description:"daemon工作目录" default-mask:"/var/run/sandbox-daemon"`
	RpcSocketPath string
}

func (o *GlobalOptions) Validate() error {
	if len(o.RuntimeDir) == 0 {
		o.RuntimeDir = "/var/run/sandbox-daemon"
	}
	if !strings.HasPrefix(o.RuntimeDir, "/") {
		return errors.New(`错误: 运行时目录必须是绝对路径 "%s"`, o.RuntimeDir)
	}
	if err := os.MkdirAll(o.RuntimeDir, 0755); err != nil {
		return errors.New(`错误: 无法创建运行时目录 "%s": %s`, o.RuntimeDir, err)
	}

	o.RpcSocketPath = o.RuntimeDir + "/rpc.sock"
	return nil
}

func (o GlobalOptions) GetRuntimeDir() string {
	return o.RuntimeDir
}
