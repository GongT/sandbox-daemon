//go:build ignore

package main

import (
	"fmt"
	"go/format"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

const magicString = "// --- GENERATED CODE --- 98accdf1-a00a-4454-a1af-65b4724ecde4"

var allTags = []string{
	"reflect",
	"config",
	"inst:deep",
	"inst",
	"proc",
	"proc:list",
	"proc:group",
	"proc:chan",
}

func main() {
	tagsFile := "./tags.go"
	bytes, err := os.ReadFile(tagsFile)
	if err != nil {
		panic(err)
	}

	content := string(bytes)
	prefix, _, exists := strings.Cut(content, magicString)

	if !exists {
		panic("magic string not found in tags.go")
	}

	str := strings.Builder{}
	str.WriteString(prefix)
	str.WriteString(magicString)
	str.WriteString("\n// DO NOT EDIT ANYTHING BELOW\n\n")

	consts := &strings.Builder{}
	methods := &strings.Builder{}

	consts.WriteString("const (\n")
	for _, tag := range allTags {
		constName := strings.ReplaceAll(tag, ":", "_")
		constName = strcase.ToScreamingSnake(constName)

		consts.WriteString("\t")
		consts.WriteString(constName)
		consts.WriteString(" DebugTag = ")
		consts.WriteString(strconv.Quote(strcase.ToSnake(tag)))
		consts.WriteString("\n")

		PascalCase := strcase.ToCamel(constName)
		fmt.Fprintf(methods, "func D%s(v ...any) {DLog(string(%s), v...)}\n", PascalCase, constName)
		fmt.Fprintf(methods, "func D%sF(fmt string, v ...any) {DLogF(string(%s), fmt, v...)}\n", PascalCase, constName)
		fmt.Fprintf(methods, "func %s(v ...any) {Log(string(%s), v...)}\n", PascalCase, constName)
		fmt.Fprintf(methods, "func %sF(fmt string, v ...any) {LogF(string(%s), fmt, v...)}\n", PascalCase, constName)
		methods.WriteString("\n")
	}
	consts.WriteString(")\n")

	str.WriteString(consts.String())
	str.WriteString("\n")
	str.WriteString(methods.String())
	str.WriteString("\n")

	bytes = []byte(str.String())
	bytes, err = format.Source(bytes)
	if err != nil {
		fmt.Println(str.String())
		panic(err)
	}

	err = os.WriteFile(tagsFile, bytes, 0644)
	if err != nil {
		panic(err)
	}
}
