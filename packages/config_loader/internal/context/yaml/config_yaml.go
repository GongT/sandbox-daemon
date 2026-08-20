package context_yaml

import (
	"reflect"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/context"
	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/paths"
	"github.com/gongt/sandbox-daemon/packages/tools/i18n/type_name"
	"gitlab.com/tozd/go/errors"
)

type configYamlContext struct {
	context.ConfigFillContext

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

func (ctx *configYamlContext) HasValue(tagPath paths.ConfigPath) (bool, error) {
	path := tagPath.JoinWithAccessor("$")
	if has, ok := ctx.eCache[path]; ok {
		return has, nil
	}

	node, exists := ctx.nCache[path]
	if !exists {
		ps, err := yaml.PathString(path)
		if err != nil {
			return false, errors.WithMessagef(err, "Yaml路径[%v]异常", path)
		}

		node, err = ps.FilterFile(ctx.astRoot)
		if err == nil {
			ctx.nCache[path] = node
		} else {
			if errors.Is(err, yaml.ErrNotFoundNode) || errors.Is(err, yaml.ErrInvalidQuery) {
				// ErrNotFoundNode: 找不到这个路径
				// ErrInvalidQuery: 路径上有object为null - 这个限定可能过于宽泛，不知道会不会漏错误
				// node = nil
			} else {
				return false, errors.WithMessagef(err, "Yaml查询[%v]时出现错误", path)
			}
		}
	}

	exists = node != nil
	ctx.eCache[path] = exists
	return exists, nil
}

func (ctx *configYamlContext) GetArraySize(tagPath paths.ConfigPath) (int, error) {
	key := tagPath.JoinWithAccessor("$")

	if node, exists := ctx.nCache[key]; !exists {
		return 0, errors.New("错误调用GetArraySize: Yaml路径不存在")
	} else {
		switch node := node.(type) {
		case *ast.SequenceNode:
			return len(node.Values), nil
		default:
			return 0, errors.New("错误调用GetArraySize: Yaml路径不是数组")
		}
	}
}

func (ctx *configYamlContext) GetObjectKeys(tagPath paths.ConfigPath) ([]string, error) {
	key := tagPath.JoinWithAccessor("$")

	if node, exists := ctx.nCache[key]; !exists {
		return nil, errors.New("错误调用GetObjectKeys: Yaml路径不存在")
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
			return nil, errors.New("错误调用GetObjectKeys: Yaml路径不是对象")
		}
	}
}

func (ctx *configYamlContext) GetValue(t reflect.Type, tagPath paths.ConfigPath) (any, error) {
	key := tagPath.JoinWithAccessor("$")

	node, exists := ctx.nCache[key]
	if !exists {
		return "", errors.New("错误调用GetValue: Yaml路径不存在")
	}

	switch node := node.(type) {
	case *ast.StringNode:
		return node.Value, nil
	case *ast.BoolNode:
		return node.Value, nil
	case *ast.IntegerNode:
		return node.Value, nil
	case *ast.FloatNode:
		return node.Value, nil
	default:
		return nil, errors.Errorf("Yaml不支持将%s转换为%s", node.Type().String(), type_name.TranslateType(t))
	}
}
