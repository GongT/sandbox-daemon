package tools

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type EnvironmentMap map[string]string

func (env *EnvironmentMap) Clear() {
	*env = EnvironmentMap{}
}

func (env *EnvironmentMap) SetLine(line string) {
	env.SetLineWithOverride(line, true)
}

func (env *EnvironmentMap) SetLineWithOverride(line string, override bool) {
	if !specialValidateLine(line) {
		return
	}

	parts := strings.SplitN(line, "=", 2)
	name := parts[0]

	if len(parts) == 2 {
		value := parts[1]
		if override {
			(*env)[name] = value
		} else {
			if _, exists := (*env)[name]; !exists {
				(*env)[name] = value
			}
		}
	} else {
		log.Printf("EnvironmentMap: 无效的环境变量行: %s", line)
	}
}

func (env *EnvironmentMap) ExtendLines(lines []string, override bool) {
	for _, line := range lines {
		env.SetLineWithOverride(line, override)
	}
}

func specialValidateLine(line string) bool {
	if strings.HasPrefix(line, "BASH_FUNC_") {
		return false
	}
	return true
}

func (env *EnvironmentMap) ToLines() []string {
	lines := make([]string, 0, len(*env))
	for key, value := range *env {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	return lines
}

type Regexp struct {
	*regexp.Regexp
}

type Set[T comparable] map[T]bool

func (s Set[T]) Has(item T) bool {
	_, exists := s[item]
	return exists
}

func loadArray[T any](node any, ptr *[]T) error {
	return nil
}
