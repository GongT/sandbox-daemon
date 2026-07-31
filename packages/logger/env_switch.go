package logger

import (
	"log"
	"os"
	"strconv"

	"github.com/gongt/sandbox-daemon/packages/myenv"
)

func init() {
	var enable bool
	if myenv.IsTesting {
		// enable = false
	} else {
		if value, ok := os.LookupEnv("LOG_TIME"); !ok {
			// enable = false
		} else {
			boolVal, err := strconv.ParseBool(value)
			if err != nil { // if the value is not a valid boolean, treat it as true if it's not empty
				enable = value != ""
			} else {
				enable = boolVal
			}
		}
	}
	if !enable {
		deleteFlag()
	}
}

// delete the date and time flags from the logger
func deleteFlag() {
	want := log.Flags() &^ (log.Ldate | log.Ltime)
	log.SetFlags(want)
}
