package reflection

import (
	"reflect"

	"github.com/gongt/sandbox-daemon/packages/logger"
	"github.com/pkg/errors"
)

// 检查t是否是标量，即布尔值、整数、浮点数、复数、字符串
func IsTypeScalar(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

// 检查v是否是一个标量值，即布尔值、整数、浮点数、复数、字符串
func IsInterfaceScalar(v any) bool {
	switch v.(type) {
	case bool, string, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64,
		complex64, complex128:
		return true
	default:
		return false
	}
}

// 检查t是否有Elem()，即指针、切片、数组、映射、通道
func IsTypeElem(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return true
	default:
		return false
	}
}

// IsContainer = IsElem && IsSerializable
func IsTypeContainer(t reflect.Type) bool {
	return IsTypeElem(t) && IsTypeSerializable(t)
}

// 检查t是否是一个可序列化的类型，即不是指针、通道、函数、接口、指针
// 这些类型无论如何都不可能序列化
// 但此函数返回true的类型不一定是完全可序列化的，比如有method的结构体
func IsTypeSerializable(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		return false
	default:
		return true
	}
}

// 返回指针指向的数据类型
func IndirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// 返回指针指向的数据值
func IndirectValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}

// 创建一个新的空值，与EmptyValue不同，EmptyValue是一个零值，可能是nil，而NewValue是一个新的实例
// 但显然这不会调用构造函数，所以如果是复杂类型，其内部的字段还是零值
func InstantiateType(t reflect.Type) reflect.Value {
	t = IndirectType(t)
	switch t.Kind() {
	case reflect.Struct:
		logger.DInst("创建新结构体")
		return reflect.New(t).Elem()
	case reflect.Slice:
		logger.DInst("创建新切片")
		return reflect.MakeSlice(t, 0, 0).Elem()
	case reflect.Map:
		logger.DInst("创建新映射")
		return reflect.MakeMap(t).Elem()
	default:
		logger.DInstF("创建新零值: %s", t)
		return reflect.New(t).Elem()
	}
}

// 是否是一个可写的指针变量
func IsUsablePointer(v reflect.Value) bool {
	return v.IsValid() && v.Kind() == reflect.Pointer && v.CanSet()
}

// 初始化多层指针，确保每一层指针本身都有效，不保证指向的值有效
// 任何一层是nil，最里层一定是nil
// 返回最里层指针
func InstantiatePointers(ptr reflect.Value) (reflect.Value, error) {
	if !ptr.IsValid() || ptr.Kind() != reflect.Pointer {
		return reflect.Value{}, errors.Errorf("参数必须是指针，例如InstantiatePointers(reflect.ValueOf(&x))而非InstantiatePointers(reflect.ValueOf(x))")
	}
	if ptr.IsNil() || !ptr.Elem().CanSet() {
		return reflect.Value{}, errors.Errorf("参数必须是可写的指针")
	}

	current := ptr
	for {
		elem := current.Elem()

		// current 已是最里层指针（其 Elem 不是指针）
		if elem.Kind() != reflect.Pointer {
			return current, nil
		}

		// elem 是下一层指针：
		// - 若它不是最里层指针且为 nil，则初始化它；
		// - 若它是最里层指针且为 nil，则保持 nil 并返回它。
		if elem.IsNil() {
			if elem.Type().Elem().Kind() == reflect.Pointer {
				elem.Set(reflect.New(elem.Type().Elem()))
			} else {
				return elem, nil
			}
		}

		current = elem
	}
}

// 检查v是否是一个需要初始化才能用的零值
func IsReadable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Chan, reflect.Slice, reflect.Array, reflect.Interface, reflect.Func:
		return !v.IsNil()
	default:
		return false
	}
}

// 检查v是否是一个能够“默认初始化”的类型
func IsDefaultInitializeMeaningful(v reflect.Type) bool {
	switch v.Kind() {
	case reflect.Map, reflect.Chan, reflect.Slice, reflect.Array, reflect.Struct:
		return true
	default:
		return false
	}
}
