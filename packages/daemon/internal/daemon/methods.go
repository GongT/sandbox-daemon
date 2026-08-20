package instance

import (
	"os"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process/mp_config"
	"github.com/gongt/sandbox-daemon/packages/rpc/rpc_server"
)

// 停止所有服务，然后重新启动守护进程
func (*D) ReExecDaemon() {
}

// 检查守护进程是否存活
func (d *D) Ping() (string, error) {
	return d.session_id, nil
}

// 启动主进程
func (*D) LaunchMainProcess(config mp_config.LifecycleConfig) error {
	return nil
}

// 停止主进程
func (*D) StopMainProcess() error {
	return nil
}

// 向主进程发送信号
func (*D) KillMainProcess(signal os.Signal) error {
	return nil
}

// 连接到主进程的标准输入/输出/错误流
func (*D) AttachStdioStreaming(ctx *rpc_server.RpcContext, stdin <-chan string, stdout chan<- string, stderr chan<- string) error {
	return nil
}
