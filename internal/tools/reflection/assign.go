package reflection

import (
	"reflect"

	"github.com/gongt/sandbox-daemon/internal/tools/i18n/type_name"
	"github.com/pkg/errors"
)

// 通过反射赋值，必须是兼容类型或可以简易转换的
// 不可用于多层指针
func AssignValue(dst any, src any) error {
	pDst := reflect.ValueOf(dst)
	if pDst.Kind() != reflect.Pointer || pDst.IsNil() {
		return errors.Errorf("目标必须是非空指针")
	}
	pSrc := reflect.Indirect(reflect.ValueOf(src))

	return AssignValueReflect(pDst.Elem(), pSrc)
}

func AssignValueReflect(dst reflect.Value, src reflect.Value) error {
	sType := src.Type()
	dType := dst.Type()

	if !dst.CanSet() {
		return errors.Errorf("目标地址不可写")
	}

	if sType.AssignableTo(dType) {
		dst.Set(src)
	} else {
		if sType.ConvertibleTo(dType) {
			convertedValue := src.Convert(dType)
			dst.Set(convertedValue)
		} else {
			return errors.Errorf("无法将%s转换为%s", type_name.TranslateType(sType), type_name.TranslateType(dType))
		}
	}
	return nil
}
