package tags

import (
	"os"

	"github.com/gongt/sandbox-daemon/packages/myenv"
)

func init() {
	if myenv.IsTesting {
		Enable("*")
	} else {
		envDebug := os.Getenv("DEBUG")
		Enable(envDebug)
	}
}
