package filesystemoperator

import "github.com/gongt/sandbox-daemon/internal/config"

type FilesystemOperator struct {
	config *config.Config
}

func New(cfg *config.Config) *FilesystemOperator {
	return &FilesystemOperator{config: cfg}
}

func (f *FilesystemOperator) Run() {
}

func (f *FilesystemOperator) Destroy() {
}
