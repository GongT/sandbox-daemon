package config

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"gitlab.com/tozd/go/errors"
)

type configYamlContext struct {
	ConfigFillContext

	text    string
	astRoot *ast.File

	nCache map[string]ast.Node
	eCache map[string]bool
}

func NewYamlContext(content string) (*configYamlContext, error) {
	astFile, err := parser.ParseBytes([]byte(content), 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &configYamlContext{
		text:    content,
		astRoot: astFile,
		nCache:  make(map[string]ast.Node),
		eCache:  make(map[string]bool),
	}, nil
}

func (ctx *configYamlContext) HasValue(tagPath ConfigPath) (bool, error) {
	path := tagPath.JoinWithAccessor("$")
	if has, ok := ctx.eCache[path]; ok {
		return has, nil
	}

	node, exists := ctx.nCache[path]
	if !exists {
		ps, err := yaml.PathString(path)
		if err != nil {
			return false, errors.WithMessagef(err, "路径[%v]异常", path)
		}

		node, err = ps.FilterFile(ctx.astRoot)
		if err == nil {
			ctx.nCache[path] = node
		} else {
			if errors.Is(err, yaml.ErrNotFoundNode) {
				// node = nil
			} else {
				return false, errors.WithStack(err)
			}
		}
	}

	exists = node != nil
	ctx.eCache[path] = exists
	return exists, nil
}

func (ctx *configYamlContext) GetArraySize(tagPath ConfigPath) (int, error) {
	key := tagPath.JoinWithAccessor("$")

	if node, exists := ctx.nCache[key]; !exists {
		return 0, errors.WithStack(fmt.Errorf("错误调用GetArraySize: 路径不存在"))
	} else {
		switch node := node.(type) {
		case *ast.SequenceNode:
			return len(node.Values), nil
		default:
			return 0, errors.WithStack(fmt.Errorf("错误调用GetArraySize: 路径不是数组"))
		}
	}
}

func (ctx *configYamlContext) GetObjectKeys(tagPath ConfigPath) ([]string, error) {
	key := tagPath.JoinWithAccessor("$")

	if node, exists := ctx.nCache[key]; !exists {
		return nil, errors.WithStack(fmt.Errorf("错误调用GetObjectKeys: 路径不存在"))
	} else {
		switch node := node.(type) {
		case *ast.MappingNode:
			keys := make([]string, len(node.Values))
			for i, v := range node.Values {
				txt := v.Key.String()
				keys[i] = txt
			}
			return keys, nil
		default:
			return nil, errors.WithStack(fmt.Errorf("错误调用GetObjectKeys: 路径不是对象"))
		}
	}
}

func (ctx *configYamlContext) GetValue(t reflect.Type, tagPath ConfigPath) (string, error) {
	key := tagPath.JoinWithAccessor("$")

	node, exists := ctx.nCache[key]
	if !exists {
		return "", errors.WithStack(fmt.Errorf("错误调用GetValue: 路径不存在"))
	}

	switch node := node.(type) {
	case *ast.StringNode:
		return node.Value, nil
	default:
		return node.String(), nil
	}
}
