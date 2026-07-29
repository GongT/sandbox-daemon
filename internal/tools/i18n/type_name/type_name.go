package type_name

import (
	"reflect"
	"strconv"
)

func PointerLevel(t reflect.Type) int {
	level := 0
	for t.Kind() == reflect.Pointer {
		level++
		t = t.Elem()
	}
	return level
}

func TranslateType(targetType reflect.Type) string {
	switch targetType.Kind() {
	case reflect.Bool:
		return "布尔值"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(targetType.Bits()) + "位整数"
	case reflect.Uint8:
		return "字节"
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.Itoa(targetType.Bits()) + "位无符号整数"
	case reflect.Float32, reflect.Float64:
		return strconv.Itoa(targetType.Bits()) + "位浮点数"
	case reflect.Complex64, reflect.Complex128:
		return strconv.Itoa(targetType.Bits()) + "位复数"
	case reflect.String:
		return "字符串"
	case reflect.Slice:
		return TranslateType(targetType.Elem()) + "数组"
	case reflect.Array:
		return "定长" + TranslateType(targetType.Elem()) + "数组(" + strconv.Itoa(targetType.Len()) + ")"
	case reflect.Map:
		return "(" + TranslateType(targetType.Key()) + "->" + TranslateType(targetType.Elem()) + ")映射"
	case reflect.Interface:
		return "任意接口"
	case reflect.Struct:
		return "结构体"
	case reflect.Uintptr:
		return "C指针"
	case reflect.UnsafePointer:
		return "任意指针"
	case reflect.Pointer:
		level := PointerLevel(targetType)
		final := indirectType(targetType)
		if level == 1 {
			return TranslateType(final) + "指针"
		}
		return strconv.Itoa(level) + "级" + TranslateType(final) + "指针"
	case reflect.Func:
		return "函数指针"
	case reflect.Chan:
		switch targetType.ChanDir() {
		case reflect.SendDir:
			return "发送" + TranslateType(targetType.Elem()) + "通道"
		case reflect.RecvDir:
			return "接收" + TranslateType(targetType.Elem()) + "通道"
		default:
			return "双向" + TranslateType(targetType.Elem()) + "通道"
		}
	default:
		return "(" + targetType.String() + ")"
	}
}

func indirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}
