package main

import (
	"github.com/gongt/sandbox-daemon/packages/rpc/generator/internal/args"
	"github.com/gongt/sandbox-daemon/packages/rpc/generator/internal/gen"
)

func main() {
	opts := args.ParseArgs()
	err := gen.DoGenerate(opts)
	if err != nil {
		panic(err)
	}
}
