package gen

import (
	"log"

	"github.com/gongt/sandbox-daemon/packages/logger"
	"github.com/gongt/sandbox-daemon/packages/rpc/generator/internal/args"
	"github.com/iancoleman/strcase"
	"gitlab.com/tozd/go/errors"
)

type generateContext struct {
	clientFile *golangFile
	serverFile *golangFile
}

func DoGenerate(opts args.Options) error {
	ctx := generateContext{
		clientFile: newOutputFile(opts.ClientOutput),
		serverFile: newOutputFile(opts.ServerOutput),
	}
	for _, file := range opts.Input {
		err := handleFile(ctx, file)
		if err != nil {
			return errors.WithMessagef(err, "解析文件 %s 失败", file)
		}
	}

	err := ctx.clientFile.Save()
	if err != nil {
		return err
	}
	err = ctx.serverFile.Save()
	if err != nil {
		return err
	}

	return nil
}

func handleFile(ctx generateContext, file string) error {
	astFile, err := readAst(file)
	if err != nil {
		return err
	}
	if astFile.Decls == nil {
		return errors.Errorf("文件 %s 没有声明", file)
	}
	logger.DLogF("包: %s\n", astFile.Name.Name)

	types := findTypes(astFile)
	if len(types) == 0 {
		return errors.Errorf("文件 %s 没有类型声明", file)
	} else if len(types) > 1 {
		return errors.Errorf("文件 %s 有多个类型声明", file)
	}

	log.Printf("类型: %s", types[0])
	outTypeName := strcase.ToLowerCamel(types[0])
	ctx.clientFile.WriteBody("type ", outTypeName, " struct {}\n\n")

	methods := methodsOf(astFile, types[0])

	for _, decl := range methods {
		log.Printf("函数: %s", decl.Name.Name)
		// godump.Dump(decl)
	}

	return nil
}
