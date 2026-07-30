package io_streams

import (
	files_handler "github.com/gongt/sandbox-daemon/internal/streams/handlers/files"
	inherit_handler "github.com/gongt/sandbox-daemon/internal/streams/handlers/inherit"
	test_handler "github.com/gongt/sandbox-daemon/internal/streams/handlers/test"
	"github.com/gongt/sandbox-daemon/internal/streams/ioconfig"
	"github.com/pkg/errors"
)

func CreateForwarder(url ioconfig.DestinationSpec) (ioconfig.LineDestination, error) {
	var r ioconfig.LineDestination
	switch url.Scheme() {
	case "inherit":
		r = inherit_handler.NewInheritForwarder()
	// case "tcp", "udp", "unix":
	// 	r = network_handler.newNetworkForwarder()
	// case "http", "https":
	// 	r = http_handler.newHttpForwarder()
	// case "ws", "wss":
	// 	r = http_handler.newWebsocketForwarder()
	case "file", "fifo", "pipe":
		r = files_handler.NewFileForwarder()
	case "test":
		r = test_handler.NewTestForwarder()
		if r != nil {
			return r, nil
		}
		fallthrough
	default:
		return nil, errors.New("未知转发器类型: " + url.Scheme())
	}
	return r, nil
}
