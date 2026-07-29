package convert

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gongt/sandbox-daemon/internal/tools/i18n/type_name"
	"github.com/gongt/sandbox-daemon/internal/tools/interfaces"
	"github.com/pkg/errors"
)

func ConvertToString[T any](value T) (string, error) {
	switch v := any(value).(type) {
	case bool:
		return strconv.FormatBool(v), nil
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10), nil
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10), nil
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64), nil
	case complex64, complex128:
		return strconv.FormatComplex(reflect.ValueOf(v).Complex(), 'g', -1, 128), nil
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	case interfaces.StringerE:
		return v.String()
	default:
		switch v := any(&value).(type) {
		case fmt.Stringer:
			return v.String(), nil
		case interfaces.StringerE:
			return v.String()
		}
		return "", errors.Errorf("%q数据无对应字符串表示", type_name.TranslateType(reflect.TypeFor[T]()))
	}
}

func ConvertStringToType(value string, targetPtr reflect.Value) error {
	if targetPtr.Type().Kind() != reflect.Pointer {
		return errors.Errorf("必须传入一个指针参数")
	}

	targetValue := targetPtr.Elem()
	targetType := targetValue.Type()

	switch targetType.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		targetValue.SetBool(val)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		numVal, err := strconv.ParseInt(value, 10, targetType.Bits())
		if err != nil {
			return err
		}
		targetValue.SetInt(numVal)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		numVal, err := strconv.ParseUint(value, 10, targetType.Bits())
		if err != nil {
			return err
		}
		targetValue.SetUint(numVal)
		return nil
	case reflect.Float32, reflect.Float64:
		numVal, err := strconv.ParseFloat(value, targetType.Bits())
		if err != nil {
			return err
		}
		targetValue.SetFloat(numVal)
		return nil
	case reflect.Complex64, reflect.Complex128:
		numVal, err := strconv.ParseComplex(value, targetType.Bits())
		if err != nil {
			return err
		}
		targetValue.SetComplex(numVal)
		return nil
	case reflect.String:
		targetValue.SetString(value)
		return nil
	}

	if v, ok := targetPtr.Interface().(interfaces.StringParser); ok {
		v.FromString(value)
		return nil
	} else if v, ok := targetPtr.Interface().(interfaces.StringParserE); ok {
		return v.FromString(value)
	}

	return errors.Errorf("不能将字符串转换为%q", type_name.TranslateType(targetType))
}

func ConvertStringTo[T any](value string, target *T) error {
	return ConvertStringToType(value, reflect.ValueOf(target))
}
