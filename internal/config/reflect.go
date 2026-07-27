package config

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gongt/sandbox-daemon/internal/tools"
	"github.com/pkg/errors"
)

type ConfigFillContext interface {
	HasValue(tagPath []string) (bool, error)
	GetArraySize(tagPath []string) (int, error)
	GetObjectKeys(tagPath []string) ([]string, error)
	GetValue(t reflect.Type, tagPath []string) (string, error)
	ConvertNonScalar(get_value string, t reflect.Type) (interface{}, error)
}

type unmarshaler interface {
	FromString(string) error
}

func WalkStruct(input any, ctx ConfigFillContext) error {
	if err := shouldBePointer(input); err != nil {
		return err
	}

	value := reflect.ValueOf(input)
	if value.IsNil() {
		return errors.WithStack(fmt.Errorf("参数不能是空指针"))
	}

	if tools.IsDebug {
		tools.DebugLogOnce("调试模式，将会打印所有的getArraySize和getValue调用")
		ctx = &loggingContext{oCtx: ctx}
	}

	return walkValue(value.Elem(), nil, ctx)
}

func shouldBePointer(v interface{}) error {
	if v == nil {
		return errors.WithStack(fmt.Errorf("参数必须是指针类型, 但实际是: <nil>"))
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.WithStack(fmt.Errorf("参数必须是指针类型, 但实际是: %v", reflect.TypeOf(v).Kind()))
	}
	return nil
}

func walkValue(v reflect.Value, tagPath []string, ctx ConfigFillContext) error {
	if !v.IsValid() {
		return nil
	}

	if !v.CanSet() {
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			return walkValue(v.Elem(), tagPath, ctx)
		}
		return nil
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !canInitializePtr(v.Type().Elem()) {
				return nil
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		return walkValue(v.Elem(), tagPath, ctx)
	}

	if v.Kind() == reflect.Interface && !v.IsNil() {
		return walkValue(v.Elem(), tagPath, ctx)
	}

	wrapErr := func(err error) error {
		if err == nil {
			return nil
		}
		errPath := strings.Join(tagPath, ".")
		if errPath == "" {
			errPath = "<root>"
		}
		return errors.Wrapf(err, "处理配置%q时出错", errPath)
	}

	if len(tagPath) > 0 {
		// 非root节点，检查是否有值，没有则直接离开
		has, err := ctx.HasValue(tagPath)
		if err != nil {
			return wrapErr(err)
		}
		if !has {
			return nil
		}
	}

	if !v.CanSet() {
		return wrapErr(errors.WithStack(fmt.Errorf("反射地址不可写")))
	}

	var err error

	// 检查v的类型是否具有自定义的转换器 (FromString) 方法
	var converter parser
	if v.CanAddr() {
		vPtr := v.Addr().Interface()
		if unmarshaler, ok := vPtr.(unmarshaler); ok {
			log.Printf("[walk] 自带FromString的类型: %s", strings.Join(tagPath, "."))
			converter = func(get_value string, t reflect.Type) (any, error) {
				err := unmarshaler.FromString(get_value)
				return nil, err
			}
			err = applyPrimitive(v, tagPath, ctx, converter)
			return wrapErr(err)
		}
	}

	switch v.Kind() {
	case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		log.Printf("[walk] 标量: %s", strings.Join(tagPath, "."))
		err = applyPrimitive(v, tagPath, ctx, nil)
	case reflect.Struct:
		log.Printf("[walk] 结构体: %s", strings.Join(tagPath, "."))
		err = walkStruct(v, tagPath, ctx)
	case reflect.Slice, reflect.Array:
		log.Printf("[walk] 数组/切片: %s", strings.Join(tagPath, "."))
		err = walkSlice(v, tagPath, ctx)
	case reflect.Map:
		log.Printf("[walk] Map: %s", strings.Join(tagPath, "."))
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
		err = walkMap(v, tagPath, ctx)
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		return errors.WithStack(fmt.Errorf("无法处理%s类型", tools.TranslateType(v.Type())))
	default:
		log.Printf("[walk] 其他类型: %s", strings.Join(tagPath, "."))
		err = applyPrimitive(v, tagPath, ctx, ctx.ConvertNonScalar)
	}

	return wrapErr(err)
}

func walkMap(v reflect.Value, tagPath []string, ctx ConfigFillContext) error {
	tMap := v.Type()
	if tMap.Key().Kind() != reflect.String {
		return errors.WithStack(fmt.Errorf("不支持map[%s]类型", tools.TranslateType(tMap.Key())))
	}

	keys, err := ctx.GetObjectKeys(tagPath)
	if err != nil {
		return err
	}

	for _, key := range keys {
		keyValue := reflect.ValueOf(key)
		elemPath := append(append([]string{}, tagPath...), key)

		// mapIndex返回的是复制
		elementValue := v.MapIndex(keyValue)
		if !elementValue.IsValid() {
			elementValue = reflect.New(tMap.Elem()).Elem()
		}

		if err := walkValue(elementValue, elemPath, ctx); err != nil {
			return err
		}

		v.SetMapIndex(keyValue, elementValue)
	}

	return nil
}

func walkStruct(v reflect.Value, tagPath []string, ctx ConfigFillContext) error {
	vType := v.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		if field.PkgPath != "" {
			continue // 跳过未导出的字段
		}
		valueField := v.Field(i)

		tag, _ := field.Tag.Lookup("config")

		var childPath []string
		if tag == "" {
			childPath = tagPath
			// 没有tag时，只有当字段是结构体时才会继续递归
			// TODO:custom class
			if valueField.Kind() != reflect.Struct {
				continue
			}
		} else {
			childPath = append(append([]string{}, tagPath...), tag)
		}
		if err := walkValue(valueField, childPath, ctx); err != nil {
			return err
		}
	}
	return nil
}

// 接收数组、切片，将配置文件中的每个元素添加到已有的数组或切片后
func walkSlice(v reflect.Value, tagPath []string, ctx ConfigFillContext) error {
	size, err := ctx.GetArraySize(append([]string{}, tagPath...))
	if err != nil {
		return err
	}
	if size == 0 {
		return nil
	}

	if size < 0 {
		return errors.WithStack(fmt.Errorf("数组大小不能为负数: %d", size))
	}

	// 如果配置文件有此项，则需要扩展切片或数组长度
	existsSize := v.Len()
	newSize := existsSize + size
	newSlice := reflect.MakeSlice(v.Type(), newSize, newSize)

	v.Set(newSlice)
	v = newSlice

	for i := range size {
		elemPath := append(append([]string{}, tagPath...), strconv.Itoa(i))
		if err := walkValue(v.Index(i+existsSize), elemPath, ctx); err != nil {
			return err
		}
	}
	return nil
}

type parser func(string, reflect.Type) (interface{}, error)

func applyPrimitive(v reflect.Value, tagPath []string, ctx ConfigFillContext, parser parser) error {
	stringRepr, err := ctx.GetValue(v.Type(), append([]string{}, tagPath...))
	if err != nil {
		return err
	}

	var convertedValue interface{}
	if parser == nil {
		convertedValue, err = tools.ConvertStringToType(stringRepr, v.Type())
	} else {
		convertedValue, err = parser(stringRepr, v.Type())
		if convertedValue == nil {
			// 实现Unmarshaler接口的类型可能会返回nil值，表示它已经在FromString方法中设置了自己的值
			return nil
		}
		if reflect.TypeOf(convertedValue) != v.Type() {
			panic("ConvertNonScalar返回的类型与预期不符")
		}
	}
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(convertedValue))
	return nil
}

func canInitializePtr(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		return true
	case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}
