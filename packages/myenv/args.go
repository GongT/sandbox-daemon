package myenv

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
)

var RpcSocketPath string
var RuntimeConfig Options

type Options struct {
	RuntimeDir string `long:"dir" description:"daemon工作目录"`
	SessionId  string `long:"sess-id" description:"(internal) 会话ID"`
}

func (o Options) GetRuntimeDir() string {
	return o.RuntimeDir
}

func (o Options) GetSessionId() string {
	return o.SessionId
}

func init() {
	_, err := flags.NewParser(&RuntimeConfig, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown).Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			switch err.Type {
			case flags.ErrHelp:
				fmt.Print(err)
				os.Exit(0)
			default:
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	// 设置运行时目录
	if RuntimeConfig.SessionId == "" {
		id, _ := uuid.NewRandom()
		RuntimeConfig.SessionId = id.String()
	}
	if RuntimeConfig.RuntimeDir == "" {
		RuntimeConfig.RuntimeDir = "/var/run/sandbox-daemon"
	}

	if !strings.HasPrefix(RuntimeConfig.RuntimeDir, "/") {
		fmt.Fprintf(os.Stderr, `错误: 运行时目录必须是绝对路径 "%s"\n`, RuntimeConfig.RuntimeDir)
	}

	if err := os.MkdirAll(RuntimeConfig.RuntimeDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, `错误: 无法创建运行时目录 "%s": %s\n`, RuntimeConfig.RuntimeDir, err)
		os.Exit(1)
	}

	RpcSocketPath = RuntimeConfig.RuntimeDir + "/rpc.sock"
}
