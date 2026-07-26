package channel

import (
	"io"
	"log"
	"os/exec"
	"time"
)

// 进程标准输入流写入器
type inputHandler struct {
	stdin io.WriteCloser
	ch    chan string

	writers *stream_set[*InputWriter]
}

func AttachInputSender(cmd *exec.Cmd) *inputHandler {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	r := &inputHandler{
		stdin:   stdin,
		ch:      make(chan string),
		writers: newStreamSet[*InputWriter](),
	}

	go func() {
		// 主转发协程，负责将ch中的数据写入stdin
		defer r.stdin.Close()

		for line := range r.ch {
			_, err := r.stdin.Write([]byte(line))
			if err != nil {
				break
			}
		}

		// 即使stdin关闭，也要继续读取ch中的数据，不能让写入方阻塞
		for range r.ch {
		}
	}()

	return r
}

// 即使程序已经退出，也可以继续GetWriter并写入数据，只是没有任何作用
func (handle *inputHandler) GetWriter() *InputWriter {
	return handle.writers.Add(newInputWriter(handle))
}

// 最终销毁，此后不能再GetWriter
func (handle *inputHandler) Destroy() {
	handle.writers.Destroy()
	close(handle.ch)
}

type InputWriter struct {
	_passive_closable

	parent *inputHandler
	ch     chan<- string
	closed <-chan struct{}
}

func newInputWriter(parent *inputHandler) *InputWriter {
	return &InputWriter{
		parent: parent,
		ch:     parent.ch,
		closed: nil,
	}
}

// 写入一行字符串，末尾会自动添加换行符，系统缓冲区满了会阻塞
func (writer *InputWriter) WriteLine(line string) bool {
	if writer.ch == nil {
		panic("InputWriter: 试图使用已经被销毁的对象")
	}
	if writer.closed == nil {
		panic("InputWriter: 初始化路径异常")
	}
	select {
	case writer.ch <- line + "\n":
		return true
	case <-writer.closed:
		return false
	case <-time.After(5 * time.Second):
		log.Println("InputWriter: 写入超时，程序的输入缓冲区可能已满")
		return false
	}
}

func (writer *InputWriter) Close() bool {
	return writer.parent.writers.Remove(writer)
}

func (writer *InputWriter) _initialize_async(abort <-chan struct{}) {
	writer.closed = abort

	go func() {
		<-abort
		writer.closed = nil
		writer.ch = nil
	}()
}
