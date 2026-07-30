package syscall_helpers

import (
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

func SetFileFlag(fd int, flag int, value bool) error {
	flags, err := unix.FcntlInt(uintptr(fd), syscall.F_GETFL, 0)
	if err != nil {
		return errors.Wrapf(err, "获取文件描述符%d的标志失败", fd)
	}

	if value {
		if flags&flag != 0 { // 已经设置了
			return nil
		}
		if _, err := unix.FcntlInt(uintptr(fd), syscall.F_SETFL, flags|flag); err != nil {
			return errors.Wrapf(err, "设置文件描述符%d的标志%d失败", fd, flag)
		}
	} else {
		if flags&flag == 0 { // 已经清除了
			return nil
		}
		if _, err := unix.FcntlInt(uintptr(fd), syscall.F_SETFL, flags&^flag); err != nil {
			return errors.Wrapf(err, "清除文件描述符%d的标志%d失败", fd, flag)
		}
	}

	return nil
}

func SetNonBlocking(fd int, value bool) error {
	return SetFileFlag(fd, syscall.O_NONBLOCK, value)
}
