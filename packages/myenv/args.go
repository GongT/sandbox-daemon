package myenv

import (
	"os"
	"strings"

	"gitlab.com/tozd/go/errors"
)

var RuntimeConfig *GlobalOptions

type GlobalOptions struct {
	RuntimeDir    string `long:"dir" description:"daemon工作目录" default:"/var/run/sandbox-daemon"`
	RpcSocketPath string
}

func (o *GlobalOptions) Validate() error {
	if !strings.HasPrefix(o.RuntimeDir, "/") {
		return errors.Errorf(`错误: 运行时目录必须是绝对路径 "%s"`, o.RuntimeDir)
	}

	if err := os.MkdirAll(o.RuntimeDir, 0755); err != nil {
		return errors.Errorf(`错误: 无法创建运行时目录 "%s": %s`, o.RuntimeDir, err)
	}

	o.RpcSocketPath = o.RuntimeDir + "/rpc.sock"
	return nil
}

func (o GlobalOptions) GetRuntimeDir() string {
	return o.RuntimeDir
}
