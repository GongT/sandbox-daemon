package config

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
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

func (c *EnvironmentMap) Unmarshal(node *yaml.Node) error {
	if c == nil || node == nil {
		log.Fatalln("错误调用 EnvironmentMap.Unmarshal: this 或 node 为 nil")
	}
	if *c == nil {
		*c = EnvironmentMap{}
	}

	if node.Kind == yaml.MappingNode {
		return node.Decode(*c)
	}

	arr := []string{}

	if err := LoadArray(node, &arr); err != nil {
		return err
	}

	c.ExtendLines(arr, true)

	return nil
}

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) Unmarshal(node *yaml.Node) error {
	if r == nil || node == nil {
		log.Fatalln("错误调用 Regexp.Unmarshal: this 或 node 为 nil")
	}
	var value string
	if err := node.Decode(&value); err != nil {
		return err
	}
	compiled, err := regexp.Compile(value)
	if err != nil {
		return err
	}
	r.Regexp = compiled
	return nil
}

type Set[T comparable] map[T]bool

func (s *Set[T]) Unmarshal(node *yaml.Node) error {
	if s == nil || node == nil {
		log.Fatalln("错误调用 Set.Unmarshal: this 或 node 为 nil")
	}

	items := []T{}
	if err := LoadArray(node, &items); err != nil {
		return err
	}

	if *s == nil {
		*s = Set[T]{}
	}

	for _, item := range items {
		(*s)[item] = true
	}

	return nil
}

func (s Set[T]) Has(item T) bool {
	_, exists := s[item]
	return exists
}
