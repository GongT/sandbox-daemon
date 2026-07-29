package deep_init

import (
	"reflect"
	"strings"

	"github.com/gongt/sandbox-daemon/internal/tools/reflection"
)

// 深度初始化一个值
// 返回所有走过的节点的指针（不仅是新添加的）
func DeepInitialize(v any) (results []any) {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Pointer || val.IsNil() || !val.Elem().CanSet() {
		panic("DeepInitialize需要一个非空指针")
	}

	var walk func(v reflect.Value, level int)
	walk = func(v reflect.Value, level int) {
		if !v.CanSet() {
			debug("%s地址不可写", strings.Repeat("  ", level))
			return
		}

		if v.Kind() == reflect.Pointer {
			var err error
			v, err = reflection.InstantiatePointers(v.Addr())
			if err != nil {
				debug("%s初始化指针失败: %v", strings.Repeat("  ", level), err)
				return
			}
			if v.IsNil() {
				debug("%s指针为nil，创建空值并指向它", strings.Repeat("  ", level))
				newVal := reflection.InstantiateType(v.Type())
				v.Set(newVal.Addr())
			}
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			debug("%s结构体", strings.Repeat("  ", level))
			// if v.CanSet(){}

			for name, field := range v.Fields() {
				debug("%s- 字段: %s (%s)", strings.Repeat("  ", level), name.Name, field.Type())
				if field.CanSet() {
					walk(field, level+1)
				} else {
					debug("%s  地址不可写", strings.Repeat("  ", level))
				}
			}
		case reflect.Slice, reflect.Array:
			debug("%s数组/切片", strings.Repeat("  ", level))
			if v.IsNil() {
				debug("%s - 初始化", strings.Repeat("  ", level))
				v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			} else {
				debug("%s - 已有", strings.Repeat("  ", level))
			}
		case reflect.Map:
			debug("%s映射", strings.Repeat("  ", level))
			if v.IsNil() {
				debug("%s - 初始化", strings.Repeat("  ", level))
				v.Set(reflect.MakeMap(v.Type()))
			} else {
				debug("%s - 已有", strings.Repeat("  ", level))
			}
			/* 通道隐式初始化问题太多了，暂不支持 */
			// case reflect.Chan:
		// 	debug("%s通道", strings.Repeat("  ", level))
		// 	if v.IsNil() {
		// 		debug("%s - 初始化", strings.Repeat("  ", level))
		// 		v.Set(reflect.MakeChan(v.Type(), 0))
		// 	} else {
		// 		debug("%s - 已有", strings.Repeat("  ", level))
		// 	}
		default:
			debug("%s跳过类型: %s", strings.Repeat("  ", level), v.Type())
		}

		results = append(results, v.Addr().Interface())
	}

	walk(val.Elem(), 0)
	return
}
