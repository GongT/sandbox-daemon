package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func LoadConfigFile(path string, configStructs ...interface{}) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadConfigObject(content, configStructs...)
}

func LoadConfigObject(object []byte, configStructs ...interface{}) ([]string, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(object, &root); err != nil {
		return nil, err
	}

	rootNode := &root
	if rootNode.Kind == yaml.DocumentNode && len(rootNode.Content) > 0 {
		rootNode = rootNode.Content[0]
	}

	unknownFields := []string{}
	for _, result := range configStructs {
		value := reflect.ValueOf(result)
		if !value.IsValid() {
			continue
		}
		if value.Kind() != reflect.Pointer || value.IsNil() {
			return nil, errors.WithStack(fmt.Errorf("config target must be a non-nil pointer"))
		}
		if err := loadConfigNode(rootNode, value.Elem(), "", &unknownFields); err != nil {
			return nil, err
		}
	}

	return unknownFields, nil
}

func loadConfigNode(node *yaml.Node, target reflect.Value, path string, unknownFields *[]string) error {
	if !target.IsValid() {
		return nil
	}

	target = makeSettableValue(target)
	if target.Kind() == reflect.Struct {
		if target.CanAddr() && target.Addr().CanInterface() {
			if _, ok := target.Addr().Interface().(interface{ Unmarshal(*yaml.Node) error }); ok {
				if method, ok := getCustomUnmarshaler(target); ok {
					return method(node)
				}
			}
		}
	}
	if !target.IsValid() {
		return nil
	}

	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			target.Set(reflect.New(target.Type().Elem()))
		}
		return loadConfigNode(node, target.Elem(), path, unknownFields)
	}

	if target.CanAddr() {
		if method, ok := getCustomUnmarshaler(target); ok {
			return method(node)
		}
	}

	switch target.Kind() {
	case reflect.Struct:
		return loadStructNode(node, target, path, unknownFields)
	case reflect.Slice:
		return loadSliceNode(node, target, unknownFields)
	case reflect.Map:
		return loadMapNode(node, target, unknownFields)
	default:
		return decodeNodeValue(node, target)
	}
}

func loadStructNode(node *yaml.Node, target reflect.Value, path string, unknownFields *[]string) error {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	consumed := map[string]bool{}
	for i := 0; i < target.NumField(); i++ {
		field := target.Type().Field(i)
		fieldValue := target.Field(i)
		fieldValue = makeSettableValue(fieldValue)
		if !fieldValue.CanSet() {
			continue
		}

		configPath := field.Tag.Get("config")
		if configPath == "" {
			if field.Type.Kind() == reflect.Struct {
				if err := loadConfigNode(node, fieldValue, path, unknownFields); err != nil {
					return err
				}
			}
			continue
		}

		segments := strings.Split(configPath, ".")
		if len(segments) > 0 {
			consumed[segments[0]] = true
		}

		resolvedNode, err := resolveNodePath(node, configPath)
		if err != nil || resolvedNode == nil {
			continue
		}
		if err := loadConfigNode(resolvedNode, fieldValue, configPath, unknownFields); err != nil {
			return err
		}
	}

	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if consumed[key] {
			continue
		}
		fullPath := key
		if path != "" {
			fullPath = path + "." + key
		}
		*unknownFields = append(*unknownFields, fullPath)
	}
	return nil
}

func makeSettableValue(value reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}
	if value.CanSet() {
		return value
	}
	if value.CanAddr() {
		return reflect.NewAt(value.Type(), unsafe.Pointer(value.Addr().Pointer())).Elem()
	}
	return value
}

func decodeNodeValue(node *yaml.Node, target reflect.Value) error {
	if node == nil {
		return nil
	}
	if !target.CanSet() && !target.CanAddr() {
		return nil
	}

	if target.Kind() == reflect.Bool {
		if node.Kind == yaml.ScalarNode {
			if node.Value == "1" || strings.EqualFold(node.Value, "true") {
				target.SetBool(true)
				return nil
			}
			if node.Value == "0" || strings.EqualFold(node.Value, "false") {
				target.SetBool(false)
				return nil
			}
		}
	}

	if target.Kind() == reflect.Int && node.Kind == yaml.ScalarNode {
		if parsed, err := strconv.Atoi(node.Value); err == nil {
			target.SetInt(int64(parsed))
			return nil
		}
		if strings.HasPrefix(node.Value, "\"") && strings.HasSuffix(node.Value, "\"") {
			if parsed, err := strconv.Atoi(strings.Trim(node.Value, `"`)); err == nil {
				target.SetInt(int64(parsed))
				return nil
			}
		}
	}

	if target.Kind() == reflect.String && node.Kind == yaml.ScalarNode {
		target.SetString(node.Value)
		return nil
	}

	if target.CanAddr() && target.CanInterface() {
		if target.Kind() == reflect.Slice || target.Kind() == reflect.Map {
			return node.Decode(target.Addr().Interface())
		}
	}

	decoded := reflect.New(target.Type())
	if err := node.Decode(decoded.Interface()); err != nil {
		return err
	}
	if target.CanSet() {
		target.Set(decoded.Elem())
	}
	return nil
}

func loadSliceNode(node *yaml.Node, target reflect.Value, unknownFields *[]string) error {
	if node == nil {
		return nil
	}
	if target.Kind() != reflect.Slice {
		return nil
	}
	if node.Kind == yaml.ScalarNode {
		if target.IsNil() {
			target.Set(reflect.MakeSlice(target.Type(), 0, 1))
		}
		elem := reflect.New(target.Type().Elem())
		if err := decodeNodeValue(node, elem.Elem()); err != nil {
			return err
		}
		target.Set(reflect.Append(target, elem.Elem()))
		return nil
	}
	if node.Kind == yaml.SequenceNode {
		if target.IsNil() {
			target.Set(reflect.MakeSlice(target.Type(), 0, len(node.Content)))
		}
		for _, child := range node.Content {
			elem := reflect.New(target.Type().Elem())
			if err := decodeNodeValue(child, elem.Elem()); err != nil {
				return err
			}
			target.Set(reflect.Append(target, elem.Elem()))
		}
		return nil
	}
	return nil
}

func loadMapNode(node *yaml.Node, target reflect.Value, unknownFields *[]string) error {
	if node == nil {
		return nil
	}
	if target.Kind() != reflect.Map {
		return nil
	}
	if target.IsNil() {
		target.Set(reflect.MakeMap(target.Type()))
	}
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := reflect.New(target.Type().Elem())
		if err := decodeNodeValue(node.Content[i+1], value.Elem()); err != nil {
			return err
		}
		target.SetMapIndex(reflect.ValueOf(key), value.Elem())
	}
	return nil
}

func resolveNodePath(node *yaml.Node, path string) (*yaml.Node, error) {
	if node == nil {
		return nil, nil
	}
	current := node
	if current.Kind == yaml.DocumentNode && len(current.Content) > 0 {
		current = current.Content[0]
	}
	if path == "" {
		return current, nil
	}

	segments := strings.Split(path, ".")
	for _, segment := range segments {
		if current.Kind != yaml.MappingNode {
			return nil, errors.WithStack(fmt.Errorf("路径 %q 无法解析 (%v)", path, current.Kind))
		}
		child := findMappingValue(current, segment)
		if child == nil {
			return nil, nil
		}
		current = child
	}
	return current, nil
}

func findMappingValue(node *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func getCustomUnmarshaler(target reflect.Value) (func(*yaml.Node) error, bool) {
	if !target.CanAddr() {
		return nil, false
	}

	method := target.Addr().MethodByName("Unmarshal")
	if !method.IsValid() {
		return nil, false
	}
	if method.Type().NumIn() != 1 || method.Type().NumOut() != 1 {
		return nil, false
	}
	if method.Type().In(0) != reflect.TypeOf((*yaml.Node)(nil)) {
		return nil, false
	}
	if method.Type().Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, false
	}

	return func(node *yaml.Node) error {
		results := method.Call([]reflect.Value{reflect.ValueOf(node)})
		if len(results) == 1 && !results[0].IsNil() {
			return results[0].Interface().(error)
		}
		return nil
	}, true
}

func LoadArray[T any](node *yaml.Node, ptr *[]T) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var firstElement T
		if err := node.Decode(&firstElement); err != nil {
			return err
		}
		*ptr = append(*ptr, firstElement)
		return nil
	case yaml.SequenceNode:
		return node.Decode(ptr)
	default:
		return errors.WithStack(fmt.Errorf("未知的节点类型: %v", node.Kind))
	}
}
