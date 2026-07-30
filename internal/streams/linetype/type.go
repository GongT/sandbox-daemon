package linetype

type LineType int

const (
	LineStdout LineType = iota
	LineStderr
	LineUnknown
)
