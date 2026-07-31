package tags

import (
	"log"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/gongt/sandbox-daemon/packages/tools/string_helper"
)

type DebugTag = string

var enabledTags = make(map[DebugTag]bool)
var debugRegStr = make(map[string]bool)
var cached_re *regexp.Regexp

func clearCache() {
	cached_re = nil
	enabledTags = make(map[DebugTag]bool)
}

func CheckEnabled(tag DebugTag) bool {
	if cached, ok := enabledTags[tag]; ok {
		return cached
	}

	v := checkRegex(tag)
	enabledTags[tag] = v
	return v
}

func checkRegex(t string) bool {
	if cached_re == nil {
		reStr := strings.Join(slices.Collect(maps.Keys(debugRegStr)), "|")
		re, err := regexp.Compile(`^(` + reStr + `)$`)
		if err != nil {
			log.Printf(`无法解析正则 "%s": %v`, reStr, err)
			enabledTags[t] = false
			return false
		}

		cached_re = re
	}

	return cached_re.MatchString(t)
}

func filterElement(input string) []string {
	input = strings.TrimSpace(regexp.MustCompile(`[\s,]+`).ReplaceAllString(input, " "))

	result := make([]string, 0)
	for _, tag := range strings.Fields(input) {
		reStr := strings.Builder{}
		for value, spId := range string_helper.CutFields(tag, []string{"*"}) {
			reStr.WriteString(regexp.QuoteMeta(value))
			if spId >= 0 {
				reStr.WriteString(".*")
			}
		}
		result = append(result, reStr.String())
	}
	return result
}

func Enable(setting DebugTag) bool {
	addedAny := false
	for _, part := range filterElement(setting) {
		if _, existed := debugRegStr[part]; !existed {
			debugRegStr[part] = true
			addedAny = true
		}
	}

	if addedAny {
		clearCache()
		checkRegex("")
	}

	return addedAny
}

func Disable(setting DebugTag) bool {
	existsAny := false
	for _, part := range filterElement(setting) {
		_, existed := debugRegStr[part]
		delete(debugRegStr, part)
		if existed {
			existsAny = true
		}
	}

	if existsAny {
		clearCache()
	}

	return existsAny
}
