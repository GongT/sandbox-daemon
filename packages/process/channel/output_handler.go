package channel

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os/exec"
	"strings"
)

// 进程标准输出、错误流读取器
// 从Cmd的对应流中读取数据，并广播给所有注册的OutputReader
type outputHandler struct {
	_passive_closable

	tag          string
	target       io.ReadCloser
	readers      *stream_set[*OutputReader]
	goingToClose bool
	cmd          *exec.Cmd
}

func AttachOutputHandler(cmd *exec.Cmd) *outputHandler {
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

func AttachErrorHandler(cmd *exec.Cmd) *outputHandler {
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
	for reader := range handle.readers.list {
		// 非阻塞写入，缓冲区满了则丢弃数据
		select {
		case reader.ch <- line:
		default:
			log.Printf("发现缓冲区已满，丢弃数据")
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
