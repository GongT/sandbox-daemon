package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goforj/godump"
	"github.com/gongt/go/pkg/signals"
	"github.com/gongt/sandbox-daemon/internal/myenv"
	"github.com/gongt/sandbox-daemon/packages/daemon"
	"github.com/jessevdk/go-flags"
	"gitlab.com/tozd/go/errors"
)

type _opts_type struct {
	*myenv.GlobalOptions

	Init daemon.InitCommand `command:"init" description:"启动守护进程"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			signals.AppQuit.Fatal(r)
		} else {
			signals.AppQuit.MainFinish()
		}
	}()

	err := wrap_main()

	if err != nil {
		panic(err)
	}
}

func test4() error {
	lm := func() error {
		return errors.WithDetails(errors.Wrap(test3(), "test4"), "test4", "hello")
	}
	return lm()
}

func test3() error {
	e1 := errors.WithMessage(test2(), "test3.1")
	e2 := errors.WithMessage(test2(), "test3.2")
	e3 := errors.WithMessage(test2(), "test3.3")
	return errors.Join(e1, e2, e3)
}

var idx int

func test2() error {
	idx++
	return errors.WithDetails(test1(), "test2", strconv.Itoa(idx))
}

func test1() error {
	return errors.WithDetails(errors.New("test1"), "test1", "wow", "code", 66)
}

func wrap_main() error {
	mainOpts := _opts_type{}

	parser := flags.NewParser(&mainOpts, flags.HelpFlag|flags.PassDoubleDash)
	parser.SubcommandsOptional = true

	if _, err := parser.Parse(); err != nil {
		if err, ok := err.(*flags.Error); ok {
			switch err.Type {
			case flags.ErrHelp:
				fmt.Print(err)
				return nil
			default:
				return (err)
			}
		} else {
			return (err)
		}
	}
	var err error

	err = mainOpts.Validate()
	if err != nil {
		return errors.WithMessage(err, "参数错误")
	}

	myenv.RuntimeConfig = mainOpts.GlobalOptions

	if parser.Active == nil {
		parser.WriteHelp(os.Stderr)
		os.Stderr.WriteString("\n")
		err = fmt.Errorf("缺少命令")
	} else {
		godump.Dump(mainOpts)

		switch parser.Active.Name {
		case "init":
			err = mainOpts.Init.Run(myenv.RuntimeConfig)
		default:
			err = fmt.Errorf("未知命令: %s", parser.Active.Name)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
