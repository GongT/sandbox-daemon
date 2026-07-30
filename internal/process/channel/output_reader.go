package channel

import (
	"time"
)

// 进程输出客户端
// 通过调用Readline或从channel中读取数据来获取输出流的内容
// 每个数据代表一行
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
