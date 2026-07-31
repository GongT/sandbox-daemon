package string_helper

import (
	"iter"
	"math"
	"strings"
)

// 使用seps中的每一个分隔符，依次切割字符串s，返回一个迭代器，每次迭代返回两个字符串：切割前的部分和用到的分隔符
// 其中分隔符是传入seps中的索引，如果没有使用分隔符（即最后一项），则返回-1
// 如果字符串本身是空的，则迭代器不会返回一个空字符串，而是直接结束迭代
// 如果字符串中不包含任何seps，则迭代器只会返回一次，返回整个字符串和-1
func CutFields(s string, seps []string) iter.Seq2[string, int] {
	left := s[:]
	right := ""

	return func(yield func(string, int) bool) {
		for {
			which := -1
			nextIndex := math.MaxInt
			for used, sep := range seps {
				next := strings.Index(left, sep)
				if next < 0 {
					continue
				}
				if next < nextIndex {
					nextIndex = next
					which = used
				}
			}

			if which == -1 {
				break
			}

			sep := seps[which]
			right = left[nextIndex+len(sep):]
			if !yield(left[:nextIndex], which) {
				return
			}
			left = right
		}

		if len(left) > 0 {
			yield(left, -1)
		}
	}
}
