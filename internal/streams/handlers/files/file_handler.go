package files_handler

import (
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
	"github.com/gongt/sandbox-daemon/internal/streams/linetype"
	"github.com/pkg/errors"
)

var _ ioconfig.LineDestination = (*fileForwarder)(nil)

type fileForwarder struct {
	Path string

	mu     sync.Mutex
	stream *os.File
	ch     <-chan linetype.LineData

	fileMode os.FileMode
	isPipe   bool
	isAppend bool
	errorCh  chan error
	reopenCh chan os.Signal
}

func NewFileForwarder() *fileForwarder {
	return &fileForwarder{}
}

func (self *fileForwarder) WriteLine(line linetype.LineData) error {
	_, err := self.stream.WriteString(line.Line)
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return errors.Errorf("写入文件失败，目标已经关闭")
		}

		return err
	}
	_, err = self.stream.WriteString("\n")
	return err
}

func (self *fileForwarder) Execute(ch <-chan linetype.LineData) error {
	for {
		select {
		case err := <-self.errorCh:
			return err
		case line, ok := <-ch:
			if !ok {
				return nil
			}
			if err := self.WriteLine(line); err != nil {
				return err
			}
		}
	}
}

func (self *fileForwarder) Destroy() (err error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.stream != nil {
		err = self.stream.Close()
		self.stream = nil
	}

	close(self.errorCh)
	if self.reopenCh != nil {
		signal.Stop(self.reopenCh)
		close(self.reopenCh)
	}

	return err
}

func (self *fileForwarder) Initialize(target ioconfig.DestinationSpec) error {
	if target.Hostname() != "" {
		return errors.Errorf("错误的文件路径: %s", target.Hostname())
	}

	self.Path = target.Path()
	modeStr := target.Get("mode")
	self.fileMode = 0666
	if modeStr != "" {
		var base int
		if strings.HasPrefix(modeStr, "0") {
			modeStr = modeStr[1:]
			base = 8
		} else {
			base = 10
		}
		mode, err := strconv.ParseInt(modeStr, base, 32)
		if err != nil {
			return errors.Wrapf(err, "未知文件模式: %s", modeStr)
		}
		self.fileMode = os.FileMode(mode & 0777)
	}

	switch target.Scheme() {
	case "file":
		self.isAppend = target.GetBool("append")
		err := self.reopen()
		if err != nil {
			return err
		}

		self.registerReopenSignal()
	case "pipe", "fifo":
		self.isPipe = true
		err := self.recreatePipe()
		if err != nil {
			return err
		}
		self.stream, err = os.OpenFile(self.Path, os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *fileForwarder) registerReopenSignal() {
	self.reopenCh = make(chan os.Signal, 1)
	signal.Notify(self.reopenCh, syscall.SIGHUP)

	go func() {
		for range self.reopenCh {
			err := self.reopen()
			if err != nil {
				err = errors.WithMessagef(err, "重新打开文件失败: %s", self.Path)
				self.errorCh <- err
			}
		}
	}()
}

func (self *fileForwarder) ensureDir() error {
	dirname := path.Dir(self.Path)
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		return errors.Wrapf(err, "创建目录%s失败", dirname)
	}
	return nil
}

func (self *fileForwarder) reopen() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.stream != nil {
		self.stream.Close()
		self.stream = nil
	}

	err := self.ensureDir()
	if err != nil {
		return err
	}

	flags := os.O_WRONLY
	if self.isAppend {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	self.stream, err = os.OpenFile(self.Path, flags, self.fileMode)
	return errors.Wrapf(err, "打开文件%s失败", self.Path)
}

func (self *fileForwarder) recreatePipe() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.stream != nil {
		self.stream.Close()
		self.stream = nil
	}

	err := self.ensureDir()
	if err != nil {
		return err
	}

	err = syscall.Mkfifo(self.Path, uint32(self.fileMode))
	if err != nil && !errors.Is(err, syscall.EEXIST) {
		return errors.Wrapf(err, "创建FIFO文件%s失败", self.Path)
	}

	flags := os.O_WRONLY | syscall.O_NONBLOCK

	self.stream, err = os.OpenFile(self.Path, flags, self.fileMode)
	return errors.Wrapf(err, "打开管道%s失败", self.Path)
}
