package config

import (
	"fmt"
	"reflect"

	"github.com/gongt/sandbox-daemon/internal/tools"
	"gitlab.com/tozd/go/errors"
)

type ConfigFillContext interface {
	HasValue(tagPath ConfigPath) (bool, error)
	GetArraySize(tagPath ConfigPath) (int, error)
	GetObjectKeys(tagPath ConfigPath) ([]string, error)
	GetValue(t reflect.Type, tagPath ConfigPath) (string, error)
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
		return errors.Errorf("参数不能是空指针")
	}

	if tools.IsDebug {
		tools.DebugLogOnce("调试模式，将会打印所有的getArraySize和getValue调用")
		ctx = &loggingContext{oCtx: ctx}
	}

	_, err := walkValue(value.Elem(), newConfigPath(), ctx)

	if err != nil {
		if wp, ok := err.(*errorWithPath); ok {
			err = wp.Unwrap()
		}
		t := reflect.TypeOf(input).Elem()
		pkgName := t.PkgPath()
		typeName := t.Name()
		return errors.WithMessagef(err, "包[%s] 类型[%s]", pkgName, typeName)
	}

	return err
}

func shouldBePointer(v interface{}) error {
	if v == nil {
		return errors.Errorf("参数必须是指针类型, 但实际是: <nil>")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.Errorf("参数必须是指针类型, 但实际是: %v", reflect.TypeOf(v).Kind())
	}
	return nil
}

func walkValue(v reflect.Value, tagPath ConfigPath, ctx ConfigFillContext) (bool, error) {
	if !v.IsValid() {
		return false, errPathF(tagPath, "值无效")
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !v.CanSet() {
				return false, errPathF(tagPath, "反射地址不可写")
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		return walkValue(v.Elem(), tagPath, ctx)
	}

	if v.Kind() == reflect.Interface && !v.IsNil() {
		return walkValue(v.Elem(), tagPath, ctx)
	}

	if tagPath.Size > 0 {
		// 非root节点，检查是否有值，没有则直接离开
		has, err := ctx.HasValue(tagPath)
		if err != nil {
			return false, errPath(tagPath, err)
		}
		if !has {
			return false, nil
		}
	}

	// 此后第一个返回值都是true

	if !v.CanSet() {
		return true, errPath(tagPath, fmt.Errorf("反射地址不可写"))
	}

	var err error

	// 检查v的类型是否具有自定义的转换器 (FromString) 方法
	var converter string_parser
	if v.CanAddr() {
		vPtr := v.Addr().Interface()
		if unmarshaler, ok := vPtr.(unmarshaler); ok {
			tools.DebugLog("[walk] 自带FromString的类型: %s", s(tagPath))
			converter = func(get_value string, t reflect.Type) (any, error) {
				err := unmarshaler.FromString(get_value)
				return nil, err
			}
			err = applyPrimitive(v, tagPath, ctx, converter)
			return true, errPath(tagPath, err)
		}
	}

	switch v.Kind() {
	case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		tools.DebugLog("[walk] 标量: %s", s(tagPath))
		err = applyPrimitive(v, tagPath, ctx, nil)
	case reflect.Struct:
		tools.DebugLog("[walk] 结构体: %s", s(tagPath))
		err = walkStruct(v, tagPath, ctx)
	case reflect.Slice, reflect.Array:
		tools.DebugLog("[walk] 数组/切片: %s", s(tagPath))
		err = walkSlice(v, tagPath, ctx)
	case reflect.Map:
		tools.DebugLog("[walk] Map: %s", s(tagPath))
		err = walkMap(v, tagPath, ctx)
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		return true, errPath(tagPath, fmt.Errorf("无法处理%s类型", tools.TranslateType(v.Type())))
	default:
		tools.DebugLog("[walk] 其他类型: %s", s(tagPath))
		err = applyPrimitive(v, tagPath, ctx, ctx.ConvertNonScalar)
	}

	return true, errPath(tagPath, err)
}

func walkMap(v reflect.Value, tagPath ConfigPath, ctx ConfigFillContext) error {
	tMap := v.Type()
	if tMap.Key().Kind() != reflect.String {
		return errPathF(tagPath, "不支持map[%s]类型", tools.TranslateType(tMap.Key()))
	}

	keys, err := ctx.GetObjectKeys(tagPath)
	if err != nil {
		return errPathW(tagPath, err, "获取对象键失败")
	}

	if v.IsNil() {
		if !v.CanSet() {
			return errPathF(tagPath, "反射地址不可写")
		}
		v.Set(reflect.MakeMap(tMap))
	}

	for _, key := range keys {
		tools.DebugLog("[walk]     + %s", key)
		keyValue := reflect.ValueOf(key)
		elemPath := tagPath.WithChild(createSegmentString(key))

		newValue := reflect.New(tMap.Elem()).Elem()

		var hasValue bool
		if hasValue, err = walkValue(newValue, elemPath, ctx); err != nil {
			return err
		}

		if hasValue {
			v.SetMapIndex(keyValue, newValue)
		}
	}

	return nil
}

func walkStruct(v reflect.Value, tagPath ConfigPath, ctx ConfigFillContext) error {
	vType := v.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		valueField := v.Field(i)

		tag, has := field.Tag.Lookup("config")

		var childPath ConfigPath
		if tag == "" {
			childPath = tagPath
			// 没有tag时，只有当字段是结构体时才会继续递归
			// TODO:custom class
			if valueField.Kind() != reflect.Struct {
				continue
			}
		} else {
			childPath = tagPath.WithChild(createSegmentString(tag))
		}

		if field.IsExported() == false {
			if has {
				return errPathF(childPath, "字段%s是未导出的，但有config标签", field.Name)
			}
			continue // 跳过未导出的字段
		}

		if _, err := walkValue(valueField, childPath, ctx); err != nil {
			return err
		}
	}
	return nil
}

// 接收数组、切片，将配置文件中的每个元素添加到已有的数组或切片后
func walkSlice(v reflect.Value, tagPath ConfigPath, ctx ConfigFillContext) error {
	cSize, err := ctx.GetArraySize(tagPath)
	if err != nil {
		return errPathW(tagPath, err, "获取数组大小失败")
	}
	if cSize == 0 {
		return nil
	}

	if cSize < 0 {
		return errPathF(tagPath, "数组大小不能为负数: %d", cSize)
	}

	// 如果配置文件有此项，则需要扩展切片或数组长度
	existsSize := v.Len()
	newSize := existsSize + cSize
	tools.DebugLog("扩展切片: 原长度: %d, 增加 %d, 新长度: %d", existsSize, cSize, newSize)
	newSlice := reflect.MakeSlice(v.Type(), newSize, newSize)

	v.Set(newSlice)
	v = newSlice

	for i := range cSize {
		tools.DebugLog("  - 元素 %d", i)
		elemPath := tagPath.WithChild(createSegmentNumber(i))
		if _, err := walkValue(v.Index(i+existsSize), elemPath, ctx); err != nil {
			return err
		}
	}
	return nil
}

type string_parser func(string, reflect.Type) (interface{}, error)

func applyPrimitive(v reflect.Value, tagPath ConfigPath, ctx ConfigFillContext, parser string_parser) error {
	stringRepr, err := ctx.GetValue(v.Type(), tagPath)
	if err != nil {
		return err
	}

	vType := v.Type()

	var convertedValue interface{}
	if parser == nil {
		convertedValue, err = tools.ConvertStringToType(stringRepr, vType)
		if err != nil {
			return errPathW(tagPath, err, "默认转换器转换失败")
		}
	} else {
		convertedValue, err = parser(stringRepr, vType)
		if err != nil {
			return errPathW(tagPath, err, "自定义转换器转换失败")
		}
		if convertedValue == nil {
			// 实现Unmarshaler接口的类型可能会返回nil值，表示它已经在FromString方法中设置了自己的值
			return nil
		}
	}

	cType := reflect.TypeOf(convertedValue)
	if !cType.AssignableTo(vType) {
		if cType.ConvertibleTo(vType) {
			convertedValue = reflect.ValueOf(convertedValue).Convert(vType).Interface()
		} else {
			return errPathF(tagPath, "数据类型不符: 预期%s, 实际%s", tools.TranslateType(vType), tools.TranslateType(cType))
		}
	}

	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(convertedValue))
	return nil
}
