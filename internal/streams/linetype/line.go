package linetype

import (
	"strconv"
	"strings"
)

type LineData struct {
	// 文本数据
	Line string
	// ID（实际就是进程pid）
	Id int
	// 类型
	Type LineType
}

func New(line string, id int, typ LineType) LineData {
	return LineData{
		Line: line,
		Id:   id,
		Type: typ,
	}
}

func (l LineData) ToJson() string {
	return `{"line":` + encode(l.Line) + `,"id":` + strconv.Itoa(l.Id) + `,"type":` + strconv.Itoa(int(l.Type)) + `}`
}

func encode(s string) string {
	var b strings.Builder
	b.WriteByte('"') // JSON 字符串必须以双引号开始

	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			// 处理其他不可见的控制字符 (ASCII 0-31)
			if r < 32 {
				// 格式化为 \uXXXX 形式 fmt.Sprintf("\\u%04x", r)
				hex := strconv.FormatInt(int64(r), 16)
				b.WriteString("\\u")
				b.WriteString(strings.Repeat("0", 4-len(hex)))
				b.WriteString(hex)
			} else {
				b.WriteRune(r)
			}
		}
	}

	b.WriteByte('"') // JSON 字符串必须以双引号结束
	return b.String()
}
