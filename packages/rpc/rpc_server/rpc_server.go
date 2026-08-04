package rpc_server

import (
	"context"
	"sync"
)

type RpcContext struct {
	ctx context.Context

	// 多个rpc调用如果需要可以共享一个锁，避免并发冲突
	Mu *sync.Mutex
}
