package internal

import (
	"reflect"

	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/context"
	context_debug "github.com/gongt/sandbox-daemon/packages/config_loader/internal/context/debug"
	"github.com/gongt/sandbox-daemon/packages/config_loader/internal/paths"
	"github.com/gongt/sandbox-daemon/packages/logger"
	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/gongt/sandbox-daemon/packages/tools/i18n/type_name"
	"github.com/gongt/sandbox-daemon/packages/tools/interfaces"
	"github.com/gongt/sandbox-daemon/packages/tools/reflection"
	"github.com/gongt/sandbox-daemon/packages/tools/reflection/deep_init"
	"gitlab.com/tozd/go/errors"
)

func WalkStruct(input any, ctx context.ConfigFillContext) error {
	if err := shouldBePointer(input); err != nil {
		return err
	}

	value := reflect.ValueOf(input)
	if value.IsNil() {
		return errors.Errorf("参数不能是空指针")
	}

	if myenv.IsDebug {
		logger.DebugLogOnce("调试模式，将会打印所有的getArraySize和getValue调用")
		ctx = context_debug.NewLoggingContext(ctx)
	}

	_, err := walkValue(value, paths.NewTwoPath(), ctx)

	if err != nil {
		if wp, ok := err.(*paths.ErrorWithPath); ok {
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

func walkValue(vPtr reflect.Value, path paths.TwoPath, ctx context.ConfigFillContext) (bool, error) {
	var err error
	vPtr, err = reflection.InstantiatePointers(vPtr) // 已确认vPtr绝对是指针
	if err != nil {
		return false, path.Err(err)
	} else if vPtr.IsNil() { // 指针本身是nil
		panic("使用DeepInitialize后，walkValue收到的vPtr出现nil值: " + path.String())
	} else if !vPtr.Elem().CanSet() { // 指向只读内存
		panic("使用DeepInitialize后，walkValue收到的vPtr地址不可写: " + path.String())
	}

	hasValue, err := ctx.HasValue(path.Tags)
	ret := func(err error) (bool, error) {
		if err != nil {
			return hasValue, path.Err(err)
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
			return ret(path.ErrF("根对象无法调用自定义转换器"))
		} else if hasValue {
			stringVal, err := ctx.GetValue(reflect.TypeFor[string](), path.Tags)
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
		logger.DConfigF("[walk] 标量: %s | %s", path.String(), vPtr.Elem().Type().String())
		err = applyPrimitive(vPtr, path, ctx)
	case reflect.Struct:
		logger.DConfigF("[walk] 结构体: %s | %s", path.String(), vPtr.Elem().Type().String())
		subHas, err = walkStruct(vPtr, path, ctx)
	case reflect.Slice, reflect.Array:
		if hasValue {
			logger.DConfigF("[walk] 数组/切片: %s | %s", path.String(), vPtr.Elem().Type().String())
			subHas, err = walkSlice(vPtr, path, ctx)
		}
	case reflect.Map:
		if hasValue {
			logger.DConfigF("[walk] Map: %s | %s", path.String(), vPtr.Elem().Type().String())
			subHas, err = walkMap(vPtr, path, ctx)
		}
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		err = errors.Errorf("无法处理%s类型", type_name.TranslateType(vPtr.Type()))
	default:
		// logger.DConfigF("[walk] 其他类型: %s", s(path.tags))
		// err = applyPrimitive(vPtr, path, ctx)
		err = errors.Errorf("无法处理%s类型", type_name.TranslateType(vPtr.Type()))
	}

	hasValue = hasValue || subHas
	return ret(err)
}

func walkMap(mapPtr reflect.Value, path paths.TwoPath, ctx context.ConfigFillContext) (hasValue bool, err error) {
	tMap := mapPtr.Elem().Type()
	if tMap.Key().Kind() != reflect.String {
		return false, path.ErrF("不支持map[%s]类型", type_name.TranslateType(tMap.Key()))
	}

	keys, err := ctx.GetObjectKeys(path.Tags)
	if err != nil {
		return false, path.ErrW(err, "获取对象键失败")
	}

	for _, key := range keys {
		logger.DConfigF("[walk]     + %s", key)

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

func walkStruct(vPtr reflect.Value, path paths.TwoPath, ctx context.ConfigFillContext) (hasValue bool, err error) {
	v := vPtr.Elem()
	vType := v.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		valuePtr := v.Field(i).Addr()

		logger.DConfigF("[walk]     + %s", field.Name)

		tag, has := field.Tag.Lookup("config")

		var childPath paths.TwoPath
		var noTag bool
		if tag == "" {
			childPath = path.WithField(field.Name, tag)
			noTag = true
		} else {
			childPath = path.WithField(field.Name, tag)
		}

		if field.IsExported() == false {
			if has {
				err = childPath.ErrF("字段%s是未导出的，但有config标签", field.Name)
				break
			}
			continue // 跳过未导出的字段
		}

		if noTag && shouldGoDeeper(field.Type) == false {
			// 没有tag时，只有当字段是结构体时才会继续递归
			logger.DConfigF("[walk]     - 没有tag且非递归字段")
			continue
		} else {
			logger.DConfigF("[walk]     - 递归字段")
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
func walkSlice(vPtr reflect.Value, path paths.TwoPath, ctx context.ConfigFillContext) (hasValue bool, err error) {
	v := vPtr.Elem()

	cSize, err := ctx.GetArraySize(path.Tags)
	if err != nil {
		return false, path.ErrW(err, "获取数组大小失败")
	}
	if cSize == 0 {
		return false, nil
	}

	if cSize < 0 {
		return false, path.ErrF("数组大小不能为负数: %d", cSize)
	}

	// 如果配置文件有此项，则需要扩展切片或数组长度
	existsSize := v.Len()
	newSize := existsSize + cSize
	logger.DConfigF("扩展切片: 原长度: %d, 增加 %d, 新长度: %d", existsSize, cSize, newSize)
	newSlice := reflect.MakeSlice(v.Type(), newSize, newSize)

	v.Set(newSlice)
	v = newSlice

	for i := range cSize {
		logger.DConfigF("  - 元素 %d", i)
		elemPath := path.WithArrayElement(i+existsSize, i)
		if _, err := walkValue(v.Index(i+existsSize).Addr(), elemPath, ctx); err != nil {
			return true, err
		}
	}
	return true, nil
}

func applyPrimitive(vPtr reflect.Value, path paths.TwoPath, ctx context.ConfigFillContext) error {
	newVal, err := ctx.GetValue(vPtr.Elem().Type(), path.Tags)
	if err != nil {
		path.Err(err)
	}

	err = reflection.AssignValueReflect(vPtr.Elem(), reflect.ValueOf(newVal))
	if err != nil {
		path.Err(err)
	}
	// err = convert.ConvertStringToType(newVal, vPtr)
	return nil
}
