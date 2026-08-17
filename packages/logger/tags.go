//go:generate go run ./generate_tags.go

package logger

import (
	"github.com/gongt/sandbox-daemon/packages/logger/internal/tags"
)

type DebugTag string

func (t DebugTag) IsEnabled() bool {
	return tags.CheckEnabled(tags.DebugTag(t))
}

func (t DebugTag) Enable() bool {
	return tags.Enable(tags.DebugTag(t))
}

func (t DebugTag) Disable() bool {
	return tags.Disable(tags.DebugTag(t))
}

func IsEnabled(tag string) bool {
	return tags.CheckEnabled(tags.DebugTag(tag))
}

func Enable(setting string) bool {
	return tags.Enable(tags.DebugTag(setting))
}

func Disable(setting string) bool {
	return tags.Disable(tags.DebugTag(setting))
}

// --- GENERATED CODE --- 98accdf1-a00a-4454-a1af-65b4724ecde4
// DO NOT EDIT ANYTHING BELOW

const (
	REFLECT    DebugTag = "reflect"
	CONFIG     DebugTag = "config"
	INST_DEEP  DebugTag = "inst:deep"
	INST       DebugTag = "inst"
	PROC       DebugTag = "proc"
	PROC_LIST  DebugTag = "proc:list"
	PROC_GROUP DebugTag = "proc:group"
	PROC_CHAN  DebugTag = "proc:chan"
	
)

func DReflect(v ...any)              { DLog(string(REFLECT), v...) }
func DReflectF(fmt string, v ...any) { DLogF(string(REFLECT), fmt, v...) }
func Reflect(v ...any)               { Log(string(REFLECT), v...) }
func ReflectF(fmt string, v ...any)  { LogF(string(REFLECT), fmt, v...) }

func DConfig(v ...any)              { DLog(string(CONFIG), v...) }
func DConfigF(fmt string, v ...any) { DLogF(string(CONFIG), fmt, v...) }
func Config(v ...any)               { Log(string(CONFIG), v...) }
func ConfigF(fmt string, v ...any)  { LogF(string(CONFIG), fmt, v...) }

func DInstDeep(v ...any)              { DLog(string(INST_DEEP), v...) }
func DInstDeepF(fmt string, v ...any) { DLogF(string(INST_DEEP), fmt, v...) }
func InstDeep(v ...any)               { Log(string(INST_DEEP), v...) }
func InstDeepF(fmt string, v ...any)  { LogF(string(INST_DEEP), fmt, v...) }

func DInst(v ...any)              { DLog(string(INST), v...) }
func DInstF(fmt string, v ...any) { DLogF(string(INST), fmt, v...) }
func Inst(v ...any)               { Log(string(INST), v...) }
func InstF(fmt string, v ...any)  { LogF(string(INST), fmt, v...) }

func DProc(v ...any)              { DLog(string(PROC), v...) }
func DProcF(fmt string, v ...any) { DLogF(string(PROC), fmt, v...) }
func Proc(v ...any)               { Log(string(PROC), v...) }
func ProcF(fmt string, v ...any)  { LogF(string(PROC), fmt, v...) }

func DProcList(v ...any)              { DLog(string(PROC_LIST), v...) }
func DProcListF(fmt string, v ...any) { DLogF(string(PROC_LIST), fmt, v...) }
func ProcList(v ...any)               { Log(string(PROC_LIST), v...) }
func ProcListF(fmt string, v ...any)  { LogF(string(PROC_LIST), fmt, v...) }

func DProcGroup(v ...any)              { DLog(string(PROC_GROUP), v...) }
func DProcGroupF(fmt string, v ...any) { DLogF(string(PROC_GROUP), fmt, v...) }
func ProcGroup(v ...any)               { Log(string(PROC_GROUP), v...) }
func ProcGroupF(fmt string, v ...any)  { LogF(string(PROC_GROUP), fmt, v...) }

func DProcChan(v ...any)              { DLog(string(PROC_CHAN), v...) }
func DProcChanF(fmt string, v ...any) { DLogF(string(PROC_CHAN), fmt, v...) }
func ProcChan(v ...any)               { Log(string(PROC_CHAN), v...) }
func ProcChanF(fmt string, v ...any)  { LogF(string(PROC_CHAN), fmt, v...) }
