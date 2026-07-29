package config

import (
	"reflect"

	"github.com/gongt/sandbox-daemon/internal/tools"
	"github.com/gongt/sandbox-daemon/internal/tools/i18n/type_name"
	"github.com/gongt/sandbox-daemon/internal/tools/interfaces"
	"github.com/gongt/sandbox-daemon/internal/tools/reflection"
	"github.com/gongt/sandbox-daemon/internal/tools/reflection/deep_init"
	"gitlab.com/tozd/go/errors"
)

type ConfigFillContext interface {
	HasValue(path ConfigPath) (bool, error)
	GetArraySize(path ConfigPath) (int, error)
	GetObjectKeys(path ConfigPath) ([]string, error)
	GetValue(t reflect.Type, path ConfigPath) (any, error)
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

	_, err := walkValue(value, newTwoPath(), ctx)

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

func shouldBePointer(v any) error {
	if v == nil {
		return errors.Errorf("参数必须是指针, 但实际是 <nil>")
	}

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Pointer {
		return errors.Errorf("参数必须是指针类型, 但得到%s", type_name.TranslateType(t))
	}
	if t.Elem().Kind() == reflect.Pointer {
		return errors.Errorf("不允许多层指针参数“%v”", t.Kind())
	}
	return nil
}

func walkValue(vPtr reflect.Value, path twoPath, ctx ConfigFillContext) (bool, error) {
	var err error
	vPtr, err = reflection.InstantiatePointers(vPtr) // 已确认vPtr绝对是指针
	if err != nil {
		return false, errPath(path, err)
	} else if vPtr.IsNil() { // 指针本身是nil
		panic("使用DeepInitialize后，walkValue收到的vPtr出现nil值: " + path.String())
	} else if !vPtr.Elem().CanSet() { // 指向只读内存
		panic("使用DeepInitialize后，walkValue收到的vPtr地址不可写: " + path.String())
	}

	hasValue, err := ctx.HasValue(path.tags)
	ret := func(err error) (bool, error) {
		if err != nil {
			return hasValue, errPath(path, err)
		}
		return hasValue, nil
	}

	if err != nil {
		return ret(err)
	}
	// 检查v的类型是否具有自定义的转换器 (FromString) 方法
	string_parser := interfaces.GetStringParser(vPtr.Interface())
	if string_parser != nil {
		// 调用自定义转换器
		if path.IsRoot() {
			return ret(errPathF(path, "根对象无法调用自定义转换器"))
		} else if hasValue {
			stringVal, err := ctx.GetValue(reflect.TypeFor[string](), path.tags)
			if err != nil {
				return ret(err)
			}
			return ret(string_parser(stringVal.(string)))
		}
	}

	// 通用类型处理
	subHas := false
	switch vPtr.Elem().Kind() {
	case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		tools.DebugLog("[walk] 标量: %s | %s", path.String(), vPtr.Elem().Type().String())
		err = applyPrimitive(vPtr, path, ctx)
	case reflect.Struct:
		tools.DebugLog("[walk] 结构体: %s | %s", path.String(), vPtr.Elem().Type().String())
		subHas, err = walkStruct(vPtr, path, ctx)
	case reflect.Slice, reflect.Array:
		tools.DebugLog("[walk] 数组/切片: %s | %s", path.String(), vPtr.Elem().Type().String())
		subHas, err = walkSlice(vPtr, path, ctx)
	case reflect.Map:
		tools.DebugLog("[walk] Map: %s | %s", path.String(), vPtr.Elem().Type().String())
		subHas, err = walkMap(vPtr, path, ctx)
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		err = errors.Errorf("无法处理%s类型", type_name.TranslateType(vPtr.Type()))
	default:
		// tools.DebugLog("[walk] 其他类型: %s", s(path.tags))
		// err = applyPrimitive(vPtr, path, ctx)
		err = errors.Errorf("无法处理%s类型", type_name.TranslateType(vPtr.Type()))
	}

	hasValue = hasValue || subHas
	return ret(err)
}

func walkMap(mapPtr reflect.Value, path twoPath, ctx ConfigFillContext) (hasValue bool, err error) {
	tMap := mapPtr.Elem().Type()
	if tMap.Key().Kind() != reflect.String {
		return false, errPathF(path, "不支持map[%s]类型", type_name.TranslateType(tMap.Key()))
	}

	keys, err := ctx.GetObjectKeys(path.tags)
	if err != nil {
		return false, errPathW(path, err, "获取对象键失败")
	}

	for _, key := range keys {
		tools.DebugLog("[walk]     + %s", key)

		vKey := reflect.ValueOf(key)
		elemPath := path.WithField(key, key)

		newValuePtr := reflect.New(tMap.Elem())
		deep_init.DeepInitialize(newValuePtr.Interface())

		if hasValue, err = walkValue(newValuePtr, elemPath, ctx); err != nil {
			return
		}

		if hasValue {
			mapPtr.Elem().SetMapIndex(vKey, newValuePtr.Elem())
		} else if !mapPtr.MapIndex(vKey).IsValid() {
			mapPtr.Elem().SetMapIndex(vKey, newValuePtr.Elem())
		}
	}

	return
}

func shouldGoDeeper(vType reflect.Type) bool {
	vType = reflection.IndirectType(vType)
	switch vType.Kind() {
	case reflect.Struct:
		return true
	case reflect.Map, reflect.Slice, reflect.Array:
		return shouldGoDeeper(vType.Elem())
	default:
		return false
	}
}

func walkStruct(vPtr reflect.Value, path twoPath, ctx ConfigFillContext) (hasValue bool, err error) {
	v := vPtr.Elem()
	vType := v.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		valuePtr := v.Field(i).Addr()

		tools.DebugLog("[walk]     + %s", field.Name)

		tag, has := field.Tag.Lookup("config")

		var childPath twoPath
		var noTag bool
		if tag == "" {
			childPath = path.WithField(field.Name, tag)
			noTag = true
		} else {
			childPath = path.WithField(field.Name, tag)
		}

		if field.IsExported() == false {
			if has {
				err = errPathF(childPath, "字段%s是未导出的，但有config标签", field.Name)
				break
			}
			continue // 跳过未导出的字段
		}

		if noTag && shouldGoDeeper(field.Type) == false {
			// 没有tag时，只有当字段是结构体时才会继续递归
			tools.DebugLog("[walk]     - 没有tag且非递归字段")
			continue
		} else {
			tools.DebugLog("[walk]     - 递归字段")
		}

		var subHas bool
		subHas, err = walkValue(valuePtr, childPath, ctx)
		hasValue = hasValue || subHas

		if err != nil {
			break
		}
	}
	return
}

// 接收数组、切片，将配置文件中的每个元素添加到已有的数组或切片后
func walkSlice(vPtr reflect.Value, path twoPath, ctx ConfigFillContext) (hasValue bool, err error) {
	v := vPtr.Elem()

	cSize, err := ctx.GetArraySize(path.tags)
	if err != nil {
		return false, errPathW(path, err, "获取数组大小失败")
	}
	if cSize == 0 {
		return false, nil
	}

	if cSize < 0 {
		return false, errPathF(path, "数组大小不能为负数: %d", cSize)
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
		elemPath := path.WithArrayElement(i+existsSize, i)
		if _, err := walkValue(v.Index(i+existsSize).Addr(), elemPath, ctx); err != nil {
			return true, err
		}
	}
	return true, nil
}

func applyPrimitive(vPtr reflect.Value, path twoPath, ctx ConfigFillContext) error {
	newVal, err := ctx.GetValue(vPtr.Elem().Type(), path.tags)
	if err != nil {
		return errPath(path, err)
	}

	err = reflection.AssignValueReflect(vPtr.Elem(), reflect.ValueOf(newVal))
	if err != nil {
		return errPath(path, err)
	}
	// err = convert.ConvertStringToType(newVal, vPtr)
	return nil
}
