package main

import (
	"fmt"
	"os"

	"github.com/goforj/godump"
	"github.com/gongt/sandbox-daemon/packages/daemon"
	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/jessevdk/go-flags"
	"golang.org/x/term"
)

type _opts_type struct {
	*myenv.GlobalOptions

	Init daemon.InitCommand `command:"init" description:"启动守护进程"`
}

func main() {
	mainOpts := _opts_type{}

	parser := flags.NewParser(&mainOpts, flags.HelpFlag|flags.PassDoubleDash)
	parser.SubcommandsOptional = false // 命令必填

	if _, err := parser.Parse(); err != nil {
		if err, ok := err.(*flags.Error); ok {
			switch err.Type {
			case flags.ErrHelp:
				fmt.Print(err)
				os.Exit(0)
			default:
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	var err error

	err = mainOpts.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "参数错误: %s\n", err)
		os.Exit(1)
	}

	myenv.RuntimeConfig = mainOpts.GlobalOptions

	godump.Dump(mainOpts)

	switch parser.Active.Name {
	case "init":
		err = mainOpts.Init.Run(myenv.RuntimeConfig)
	default:
		err = fmt.Errorf("未知命令: %s", parser.Active.Name)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %s\n", err)
		os.Exit(1)
	}

	if term.IsTerminal(int(os.Stderr.Fd())) {
		fmt.Println("bye, bye.")
	}
}
