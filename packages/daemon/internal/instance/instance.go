package instance

import (
	"os/signal"
	"sync"
	"syscall"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/reap"
	"github.com/gongt/sandbox-daemon/packages/myenv"
)

type D struct {
	mu       sync.RWMutex
	quitChan chan struct{}

	sessionId string
	config    *myenv.GlobalOptions

	parts []internal.DaemonComponent
}

func New(config internal.WithSessionId, runtime *myenv.GlobalOptions) *D {
	// 防止stdout、stderr被关闭时，程序直接退出
	signal.Reset(syscall.SIGPIPE)

	i := D{
		quitChan:  make(chan struct{}),
		sessionId: config.GetSessionId(),
		config:    runtime,
	}

	// rpcServer := rpc.NewServer()
	reaper := reap.New()
	mp := main_process.New()

	i.parts = []internal.DaemonComponent{
		// rpcServer,
		reaper,
		mp,
	}

	return &i
}
