package environ

import (
	"fmt"
	"maps"
	"strings"

	"github.com/pkg/errors"
)

type Map map[string]string

func (env *Map) Clear() {
	*env = Map{}
}

func (env *Map) SetLine(line string) error {
	return env.SetLineWithOverride(line, true)
}

func (env *Map) SetLineWithOverride(line string, override bool) error {
	if !specialValidateLine(line) {
		return nil
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
		return nil
	} else {
		return errors.Errorf("EnvironmentMap: 无效的环境变量行: %s", line)
	}
}

func (env *Map) Extend(other map[string]string, override bool) {
	if override {
		maps.Copy(*env, other)
	} else {
		for key, value := range other {
			if _, exists := (*env)[key]; !exists {
				(*env)[key] = value
			}
		}
	}
}

func (env *Map) ExtendLines(lines []string, override bool) error {
	for _, line := range lines {
		err := env.SetLineWithOverride(line, override)
		if err != nil {
			return err
		}
	}
	return nil
}

func specialValidateLine(line string) bool {
	if strings.HasPrefix(line, "BASH_FUNC_") {
		return false
	}
	return true
}

func (env *Map) ToLines() []string {
	lines := make([]string, 0, len(*env))
	for key, value := range *env {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	return lines
}

func (env *Map) Delete(name string) {
	delete(*env, name)
}

func (env *Map) Has(name string) bool {
	_, exists := (*env)[name]
	return exists
}

func (env *Map) Set(name, value string) {
	(*env)[name] = value
}

func (env *Map) Get(name string) string {
	value, _ := (*env)[name]
	return value
}

func (env *Map) Size() int {
	return len(*env)
}
