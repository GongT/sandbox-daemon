package io_streams

import (
	"io"
	"sync"

	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
	"github.com/gongt/sandbox-daemon/internal/streams/linetype"
	"github.com/gongt/sandbox-daemon/internal/tools"
)

type holder struct {
	id uint64
	// 目标输出
	destination ioconfig.LineDestination
	channel     chan linetype.LineData
	done        chan struct{}

	closed sync.Once
	err    error
}

func (b *holder) Close() (err error) {
	b.closed.Do(func() { // close有两条路径，需要保护close
		close(b.channel)
		for range b.channel {
			// drain
		}
		<-b.done
		err = b.destination.Destroy()
		b.destination = nil
	})
	if err == nil {
		err = b.err
	} else if b.err == nil {
		b.err = err
	}
	return
}

type LineMultiplexer struct {
	// 目标输出
	destinations []*holder

	mu   sync.Mutex
	done chan struct{}
}

func NewMultiplexer() *LineMultiplexer {
	return &LineMultiplexer{
		destinations: make([]*holder, 0),
		done:         make(chan struct{}),
	}
}

func (m *LineMultiplexer) AddDestination(target ioconfig.DestinationSpec) (io.Closer, error) {
	dest, err := CreateForwarder(target)
	if err != nil {
		return nil, err
	}

	err = dest.Initialize(target)
	if err != nil {
		return nil, err
	}

	id := tools.GetId()
	box := &holder{
		id:          id,
		destination: dest,
		channel:     make(chan linetype.LineData, 100),
		done:        make(chan struct{}),
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.destinations = append(m.destinations, box)

	go func() {
		err = dest.Execute(box.channel)
		m.delete(id)
		box.err = err
		close(box.done)
		box.Close()
	}()

	return box, nil
}

func (m *LineMultiplexer) delete(id uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, d := range m.destinations {
		if d.id == id {
			m.destinations = append(m.destinations[:i], m.destinations[i+1:]...)
			break
		}
	}
}

func (m *LineMultiplexer) WriteLine(line linetype.LineData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, dest := range m.destinations {
		select {
		case dest.channel <- line:
		default:
		}
	}
	return nil
}
