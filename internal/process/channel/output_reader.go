package channel

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

// 进程标准输出、错误流读取器
type outputHandler struct {
	_passive_closable

	tag          string
	target       io.ReadCloser
	readers      *stream_set[*OutputReader]
	goingToClose bool
	cmd          *exec.Cmd
}

func AttachOutputReader(cmd *exec.Cmd) *outputHandler {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	r := &outputHandler{
		tag:     "stdout",
		target:  stdout,
		readers: newStreamSet[*OutputReader](),
		cmd:     cmd,
	}

	go r.execute()

	return r
}

func AttachErrorReader(cmd *exec.Cmd) *outputHandler {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	r := &outputHandler{
		tag:     "stderr",
		target:  stderr,
		readers: newStreamSet[*OutputReader](),
		cmd:     cmd,
	}

	go r.execute()

	return r
}

func (handle *outputHandler) execute() {
	bread := bufio.NewReader(handle.target)

	lbuf := make([]byte, 0, 4096)
	for {
		bytes, isFull, err := bread.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else if strings.Contains(err.Error(), "file already closed") {
				break
			}
			panic(err)
		}

		var lineText string
		if len(lbuf) == 0 && !isFull {
			// 一次性读取到一整行，直接转换为string，减少一次内存拷贝
			lineText = string(bytes)
		} else {
			// 需拼成一整行才能转换，否则可能会出现乱码
			lbuf = append(lbuf, bytes...)
			if isFull {
				continue
			}
			lineText = string(lbuf)
		}
		lbuf = lbuf[:0]

		handle.broadcast(lineText)
	}

	if !handle.goingToClose {
		log.Printf("进程%s流读取器提前关闭", handle.tag)
	}
}

func (handle *outputHandler) broadcast(line string) {
	for reader, _ := range handle.readers.list {
		// 非阻塞写入，若没有任何reader，或者缓冲区满了，则丢弃数据
		select {
		case reader.ch <- line:
		default:
			log.Printf("发现")
		}
	}
}

func (handle *outputHandler) GetReader(name string) *OutputReader {
	return handle.readers.Add(newOutputReader(name, handle))
}

func (handle *outputHandler) Destroy() {
	err := handle.target.Close()
	if err != nil {
		panic(err)
	}
	handle.readers.Destroy()
}

type OutputReader struct {
	_passive_closable

	parent       *outputHandler
	name         string
	ch           chan string
	close_signal <-chan struct{}
}

func newOutputReader(name string, parent *outputHandler) *OutputReader {
	return &OutputReader{
		parent:       parent,
		name:         name,
		ch:           make(chan string, 100),
		close_signal: nil,
	}
}

func (reader *OutputReader) Wait() <-chan struct{} {
	return reader.close_signal
}

func (reader *OutputReader) GetChannel() <-chan string {
	reader._assert_ok()
	return reader.ch
}

func (reader *OutputReader) Readline() (string, bool) {
	reader._assert_ok()

	select {
	case line := <-reader.ch:
		return line, true
	case <-reader.close_signal:
		return "", false
	}
}

func (reader *OutputReader) ReadlineTimeout(timeout time.Duration) (string, bool) {
	reader._assert_ok()

	select {
	case line := <-reader.ch:
		return line, true
	case <-reader.close_signal:
		return "", false
	case <-time.After(timeout):
		return "", false
	}
}

func (reader *OutputReader) _assert_ok() {
	if reader.ch == nil {
		panic("OutputReader: 试图使用已经被销毁的对象")
	}
	if reader.close_signal == nil {
		panic("OutputReader: 初始化路径异常")
	}
}

func (reader *OutputReader) Close() bool {
	reader.ch = nil
	return reader.parent.readers.Remove(reader)
}

func (reader *OutputReader) _initialize_async(abort <-chan struct{}) {
	reader.close_signal = abort
	go func() {
		<-abort
		close(reader.ch)
	}()
}
