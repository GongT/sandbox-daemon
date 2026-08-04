package args

import (
	"os"
	"strings"

	"github.com/gongt/sandbox-daemon/packages/rpc/generator/internal/tools"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Input        []string `long:"in" description:"输入文件列表"`
	ServerOutput string   `long:"server" description:"输出文件路径"`
	ClientOutput string   `long:"client" description:"输出文件路径"`
}

func ParseArgs() Options {
	opts := Options{}
	_, err := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown).Parse()
	if err != nil {
		panic(err)
	}

	if len(opts.Input) == 0 {
		file := os.Getenv("GOFILE")
		if file != "" {
			opts.Input = []string{file}
		}
	}
	opts.Input = tools.AbsoluteList(opts.Input)
	if len(opts.Input) == 0 {
		panic("没有输入文件")
	}

	// out 1
	if opts.ServerOutput == "" {
		file := os.Getenv("GOFILE")
		if strings.HasSuffix(file, ".go") {
			opts.ServerOutput = file[:len(file)-3] + ".server.go"
		}
	}
	opts.ServerOutput = tools.Absolute(opts.ServerOutput)
	if opts.ServerOutput == "" {
		panic("没有-server输出文件")
	}

	// out 2
	if opts.ClientOutput == "" {
		file := os.Getenv("GOFILE")
		if strings.HasSuffix(file, ".go") {
			opts.ClientOutput = file[:len(file)-3] + "_client/main.go"
		}
	}
	opts.ClientOutput = tools.Absolute(opts.ClientOutput)
	if opts.ClientOutput == "" {
		panic("没有-client输出文件")
	}

	return opts
}
