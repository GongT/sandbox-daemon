//go:generate go run ../../../rpc/generator/main.go --client ../../client/handler.generated.go --server ./handler.generated.go
package rpc

import (
	"os"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process/config"
	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/gongt/sandbox-daemon/packages/rpc/rpc_server"
)

type RpcHandler struct{}

// 停止所有服务，然后重新启动守护进程
func (*RpcHandler) ReExecDaemon() {
}

// 检查守护进程是否存活
func (*RpcHandler) Ping() (string, error) {
	return myenv.RunSessionId, nil
}

// 启动主进程
func (*RpcHandler) LaunchMainProcess(config config.LifecycleConfig) error {
	return nil
}

// 停止主进程
func (*RpcHandler) StopMainProcess() error {
	return nil
}

// 向主进程发送信号
func (*RpcHandler) KillMainProcess(signal os.Signal) error {
	return nil
}

// 连接到主进程的标准输入/输出/错误流
func (*RpcHandler) AttachStdioStreaming(ctx *rpc_server.RpcContext, stdin <-chan string, stdout chan<- string, stderr chan<- string) error {
	return nil
}
