//go:build ignore

package main

import "github.com/gongt/sandbox-daemon/packages/logger"

func main() {
	logger.DLogF("test", "wow, such %s", "doge")
}
