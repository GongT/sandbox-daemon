package rpc

import "github.com/gongt/sandbox-daemon/packages/daemon/internal"

var _ internal.DaemonComponent = (*rpcServer)(nil)

type rpcServer struct{}

func NewServer() *rpcServer {
	return &rpcServer{}
}

func (r *rpcServer) Start() error {
	return nil
}

func (r *rpcServer) Stop() error {
	return nil
}
