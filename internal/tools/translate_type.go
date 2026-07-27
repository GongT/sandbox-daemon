package tools

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

func TranslateType(targetType reflect.Type) string {
	switch targetType.Kind() {
	case reflect.Bool:
		return "布尔值"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "整数"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "无符号整数"
	case reflect.Float32, reflect.Float64:
		return "浮点数"
	case reflect.Complex64, reflect.Complex128:
		return "复数"
	case reflect.String:
		return "字符串"
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return "数组对象"
	case reflect.Map:
		return "映射对象"
	case reflect.Interface:
		return "任意接口"
	case reflect.Struct:
		return "结构体对象"
	case reflect.UnsafePointer:
		fallthrough
	case reflect.Pointer:
		fallthrough
	case reflect.Uintptr:
		return "指针"
	case reflect.Func:
		return "函数指针"
	case reflect.Chan:
		return "通道"
	default:
		return targetType.String()
	}
}

func ConvertStringToType(value string, targetType reflect.Type) (any, error) {
	switch targetType.Kind() {
	case reflect.Bool:
		return strconv.ParseBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		numVal, err := strconv.ParseInt(value, 10, targetType.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(numVal).Convert(targetType).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		numVal, err := strconv.ParseUint(value, 10, targetType.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(numVal).Convert(targetType).Interface(), nil
	case reflect.Float32, reflect.Float64:
		numVal, err := strconv.ParseFloat(value, targetType.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(numVal).Convert(targetType).Interface(), nil
	case reflect.Complex64, reflect.Complex128:
		numVal, err := strconv.ParseComplex(value, targetType.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(numVal).Convert(targetType).Interface(), nil
	case reflect.String:
		return value, nil
	default:
		return nil, errors.WithStack(fmt.Errorf("%q数据无对应字符串表示", TranslateType(targetType)))
	}
}

func ConvertTypeToString(value any) (string, error) {
	typ := reflect.TypeOf(value)
	switch typ.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(value.(bool)), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(reflect.ValueOf(value).Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(reflect.ValueOf(value).Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(reflect.ValueOf(value).Float(), 'f', -1, typ.Bits()), nil
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(reflect.ValueOf(value).Complex(), 'g', -1, typ.Bits()), nil
	case reflect.String:
		return value.(string), nil
	default:
		return "", errors.WithStack(fmt.Errorf("%q数据无对应字符串表示", TranslateType(typ)))
	}
}
