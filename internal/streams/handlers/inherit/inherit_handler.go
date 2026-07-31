package inherit_handler

import (
	"os"
	"strconv"
	"syscall"

	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
	"github.com/gongt/sandbox-daemon/internal/streams/linetype"
	"github.com/gongt/sandbox-daemon/packages/tools/syscall_helpers"
	"github.com/pkg/errors"
)

var _ ioconfig.LineDestination = (*inheritForwarder)(nil)

type inheritForwarder struct {
	stream       *os.File
	fd           int
	errorCounter int
	ch           <-chan linetype.LineData
	muteError    bool
}

func NewInheritForwarder() *inheritForwarder {
	return &inheritForwarder{}
}

func (self *inheritForwarder) WriteLine(line linetype.LineData) error {
	_, err := self.stream.WriteString(line.Line)
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return errors.Errorf("写入文件描述符%d失败，目标已经关闭", self.fd)
		}
		self.errorCounter++

		if self.errorCounter >= 10 {
			switch self.fd {
			case 1:
				return errors.WithMessage(err, "写入标准输出多次失败")
			case 2:
				return errors.WithMessage(err, "写入标准错误多次失败")
			default:
				return errors.WithMessagef(err, "写入文件描述符%d多次失败", self.fd)
			}
		}

		return err
	}
	_, err = self.stream.WriteString("\n")
	return err
}

func (self *inheritForwarder) Execute(ch <-chan linetype.LineData) error {
	for line := range ch {
		if err := self.WriteLine(line); err != nil {
			if self.muteError {
				return nil
			} else {
				return err
			}
		}
	}

	return nil
}

// 绝不能关闭描述符，因为有可能在重启路径中，关了就寄了
func (self *inheritForwarder) Destroy() error {
	self.stream.Sync()
	// 需要改回阻塞模式？
	return nil
}

func (self *inheritForwarder) Initialize(target ioconfig.DestinationSpec) error {
	switch target.Hostname() {
	case "stdout":
		self.stream = os.Stdout
		self.fd = 1
	case "stderr":
		self.stream = os.Stderr
		self.fd = 2
	default:
		if fd, err := strconv.Atoi(target.Hostname()); err == nil {
			file := os.NewFile(uintptr(fd), "fd"+target.Hostname())
			if _, err := file.Stat(); err != nil {
				return errors.Errorf("继承的文件描述符 %d 不存在", fd)
			}
			self.stream = file
			self.fd = fd
		} else {
			return errors.Errorf("不支持的继承输出目标: %s", target.Hostname())
		}
	}

	if target.Get("error") == "ignore" {
		self.muteError = true
	}

	err := syscall_helpers.SetNonBlocking(self.fd, true)
	if err != nil && !self.muteError {
		return err
	}
	return nil
}
