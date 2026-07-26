package channel

import (
	"fmt"
	"sync"
)

type _passive_closable interface {
	// 实现者须在此方法中等待abort信号，并在收到信号后关闭自身资源
	_initialize_async(abort <-chan struct{})
}

type closable_comparable interface {
	_passive_closable
	comparable
}

type stream_set[T closable_comparable] struct {
	mu   sync.Mutex
	list map[T]chan struct{}
}

func newStreamSet[T closable_comparable]() *stream_set[T] {
	return &stream_set[T]{
		list: make(map[T]chan struct{}),
	}
}

func (s *stream_set[T]) Add(stream T) T {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.assert()

	if _, exists := s.list[stream]; exists {
		panic(fmt.Sprintf("duplicate stream added to stream_set: %v", stream))
	}

	ch := make(chan struct{})
	stream._initialize_async(ch)
	s.list[stream] = ch

	return stream
}

func (s *stream_set[T]) Remove(stream T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.assert()

	if _, exists := s.list[stream]; exists {
		close(s.list[stream])
		delete(s.list, stream)
		return true
	}
	return false
}

func (s *stream_set[T]) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.assert()

	return len(s.list) == 0
}

func (s *stream_set[T]) Destroy() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.assert()

	for _, ch := range s.list {
		close(ch)
	}
	s.list = nil
}

func (s *stream_set[T]) assert() {
	if s.list == nil {
		panic("试图使用已经被销毁的stream_set")
	}
}
