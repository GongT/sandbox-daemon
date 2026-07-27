package config

import "strconv"

type segment struct {
	isNumber bool

	str_value string
	int_value int
}

func createSegmentString(value string) segment {
	return segment{
		isNumber:  false,
		str_value: value,
	}
}

func createSegmentNumber(value int) segment {
	return segment{
		isNumber:  true,
		int_value: value,
		str_value: strconv.Itoa(value),
	}
}

func (s segment) String() string {
	return s.str_value
}

func (s segment) Number() int {
	if !s.isNumber {
		panic("segment is not a number")
	}
	return s.int_value
}

func (s segment) Value() any {
	if s.isNumber {
		return s.int_value
	}
	return s.str_value
}

type ConfigPath struct {
	segments []segment
	Size     int
}

func newConfigPath() ConfigPath {
	return ConfigPath{
		segments: make([]segment, 0),
		Size:     0,
	}
}

func (p ConfigPath) JoinWithDot(prefix string) string {
	result := prefix
	for _, seg := range p.segments {
		if result != "" {
			result += "."
		}
		result += seg.String()
	}
	return result
}

func (p ConfigPath) JoinWithAccessor(prefix string) string {
	result := prefix
	for _, seg := range p.segments {
		if seg.isNumber {
			result += "[" + seg.String() + "]"
		} else {
			if result != "" {
				result += "."
			}
			result += seg.String()
		}
	}
	return result
}

func (p ConfigPath) WithChild(child segment) ConfigPath {
	if child.isNumber && len(p.segments) == 0 {
		panic("根元素不能是数组")
	}

	r := make([]segment, 0, len(p.segments)+1)
	r = append(r, p.segments...)
	r = append(r, child)

	return ConfigPath{
		segments: r,
		Size:     len(r),
	}
}

func (p ConfigPath) StringAt(index int) string {
	if index < 0 {
		index = len(p.segments) + index
	}
	if index < 0 || index >= len(p.segments) {
		panic("index out of range")
	}
	return p.segments[index].String()
}
