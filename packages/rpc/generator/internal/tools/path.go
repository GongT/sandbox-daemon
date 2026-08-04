package tools

import "path/filepath"

func AbsoluteList(files []string) []string {
	abs := make([]string, len(files))
	for i, f := range files {
		abs[i] = Absolute(f)
	}
	return abs
}

func Absolute(file string) string {
	abs, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	return abs
}
