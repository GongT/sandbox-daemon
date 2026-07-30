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

// 里面实际只有3个指针，调用方不需要*DestinationSpec
type DestinationSpec struct {
	raw   string
	url   *url.URL
	query *url.Values
}

func NewDestination(raw string) (*DestinationSpec, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	spec := DestinationSpec{
		raw:   raw,
		url:   u,
		query: &query,
	}

	return &spec, nil
}

func (s *DestinationSpec) FromString(str string) error {
	spec, err := NewDestination(str)
	if err != nil {
		return err
	}
	*s = *spec
	return nil
}

func (s *DestinationSpec) Raw() string {
	return s.raw
}

func (s *DestinationSpec) Get(key string) string {
	return s.query.Get(key)
}

func (s *DestinationSpec) Has(key string) bool {
	return s.query.Has(key)
}

func (s *DestinationSpec) GetBool(key string) bool {
	b, _ := strconv.ParseBool(s.query.Get(key))
	return b
}

// clone URL methods
func (s *DestinationSpec) Scheme() string {
	return s.url.Scheme
}

func (s *DestinationSpec) Username() string {
	if s.url.User == nil {
		return ""
	}
	return s.url.User.Username()
}

func (s *DestinationSpec) Password() string {
	if s.url.User == nil {
		return ""
	}
	password, _ := s.url.User.Password()
	return password
}

func (s *DestinationSpec) Hostname() string {
	return s.url.Hostname()
}

func (s *DestinationSpec) Port() int {
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

func (s *DestinationSpec) Domain() string {
	return s.url.Host
}

func (s *DestinationSpec) Path() string {
	return s.url.Path
}

func (s *DestinationSpec) Fragment() string {
	return s.url.Fragment
}
