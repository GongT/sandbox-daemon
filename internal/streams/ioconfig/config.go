package ioconfig

import (
	"net/url"
	"strconv"

	"github.com/gongt/sandbox-daemon/internal/streams/linetype"
)

// 文本目标接口
type LineDestination interface {
	// 初始化目标，例如创建对应文件、网络连接等
	// 只会被调用一次，且在 Attach 之前
	Initialize(target DestinationSpec) error
	// 写入数据，只会被调用一次
	// 应[同步]地读取 ch 中的数据，直到 ch 被关闭
	// 提前返回调用方随后也会关闭ch
	Execute(ch <-chan linetype.LineData) error
	// 此时可以销毁资源
	// 一定在Attach返回之后调用，且只会被调用一次
	Destroy() error
}

type DestinationSpec interface {
	Raw() string
	Get(key string) string
	Has(key string) bool
	GetBool(key string) bool

	Scheme() string
	Username() string
	Password() string
	Hostname() string
	Port() int
	Domain() string
	Path() string
	Fragment() string
}

type destinationSpec struct {
	raw   string
	url   *url.URL
	query url.Values
}

func NewDestination(raw string) (DestinationSpec, error) {
	spec := &destinationSpec{raw: raw}
	if err := spec._init(); err != nil {
		return nil, err
	}
	return spec, nil
}

func (s *destinationSpec) FromString(str string) error {
	s.raw = str
	return s._init()
}

func (s *destinationSpec) _init() (err error) {
	s.url, err = url.Parse(s.raw)
	if err != nil {
		return
	}
	if s.url.RawQuery != "" {
		s.query, err = url.ParseQuery(s.url.RawQuery)
	} else {
		s.query = make(url.Values)
	}
	return
}

func (s *destinationSpec) Raw() string {
	return s.raw
}

func (s *destinationSpec) Get(key string) string {
	return s.query.Get(key)
}

func (s *destinationSpec) Has(key string) bool {
	return s.query.Has(key)
}

func (s *destinationSpec) GetBool(key string) bool {
	b, _ := strconv.ParseBool(s.query.Get(key))
	return b
}

// clone URL methods
func (s *destinationSpec) Scheme() string {
	return s.url.Scheme
}

func (s *destinationSpec) Username() string {
	if s.url.User == nil {
		return ""
	}
	return s.url.User.Username()
}

func (s *destinationSpec) Password() string {
	if s.url.User == nil {
		return ""
	}
	password, _ := s.url.User.Password()
	return password
}

func (s *destinationSpec) Hostname() string {
	return s.url.Hostname()
}

func (s *destinationSpec) Port() int {
	p := s.url.Port()
	if p == "" {
		return -1
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return -2
	}
	return port
}

func (s *destinationSpec) Domain() string {
	return s.url.Host
}

func (s *destinationSpec) Path() string {
	return s.url.Path
}

func (s *destinationSpec) Fragment() string {
	return s.url.Fragment
}
