package internal

type DaemonComponent interface {
	// 启动组件
	Start() error

	// 停止组件
	Stop() error
}
