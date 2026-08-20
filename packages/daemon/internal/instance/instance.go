package instance

import (
	"io"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gongt/sandbox-daemon/packages/daemon/internal"
	daemon_config "github.com/gongt/sandbox-daemon/packages/daemon/internal/config"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/main_process"
	"github.com/gongt/sandbox-daemon/packages/daemon/internal/reap"
	"github.com/gongt/sandbox-daemon/packages/tools/types"
	"github.com/pkg/errors"
)

type D struct {
	mu       sync.RWMutex
	quitChan chan struct{}

	parts types.Collection[io.Closer]
	mp    *main_process.MainProcess

	sessionId string
	config    *daemon_config.DaemonConfig
}

func New(options internal.WithSessionId, config *daemon_config.DaemonConfig) *D {
	// 防止stdout、stderr被关闭时，程序直接退出
	signal.Reset(syscall.SIGPIPE)

	i := D{
		quitChan:  make(chan struct{}),
		sessionId: options.GetSessionId(),
		config:    config,
	}

	reaper := reap.New()
	i.mp = main_process.New()

	i.parts = types.Collection[io.Closer]{
		reaper,
		i.mp,
	}

	return &i
}

func (d *D) Assert() error {
	if d == nil {
		return errors.New("缺少守护进程实例")
	}
	return nil
}
