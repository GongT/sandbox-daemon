package type_name

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateType(t *testing.T) {
	log.SetOutput(t.Output())

	assert.Equal(t, "布尔值", TranslateType(reflect.TypeFor[bool]()))
	assert.Equal(t, "任意接口", TranslateType(reflect.TypeFor[any]()))

	var a *int
	assert.Equal(t, "64位整数指针", TranslateType(reflect.TypeOf(a)))

	var b *****string
	assert.Equal(t, "5级字符串指针", TranslateType(reflect.TypeOf(b)))

	var c **map[string]int
	assert.Equal(t, "2级(字符串->64位整数)映射指针", TranslateType(reflect.TypeOf(c)))

	var d chan<- *float32
	assert.Equal(t, "发送32位浮点数指针通道", TranslateType(reflect.TypeOf(d)))

	var e [6]string
	assert.Equal(t, "定长字符串数组(6)", TranslateType(reflect.TypeOf(e)))

	var f []byte
	assert.Equal(t, "字节数组", TranslateType(reflect.TypeOf(f)))
}
