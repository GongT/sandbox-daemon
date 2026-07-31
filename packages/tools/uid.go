package tools

var g uint64

// 从1开始自增
func GetId() uint64 {
	g++
	return g
}
