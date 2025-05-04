package util

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

var (
	DEBUG     = false
	ConfigDir = "."
)

type DebugInfo struct {
	ServerTags  *OrderedMap[string]
	AppTags     *OrderedMap[string]
	RuntimeTags *OrderedMap[string]
	Mods        []*debug.Module
}

func DebugGetInfo(version string, started time.Time) *DebugInfo {
	bi := DebugBuildInfo()

	serverTags := NewOrderedMap[string](false, 10)
	hostname, _ := os.Hostname()
	serverTags.Set("Machine Name", hostname)
	serverTags.Set("CPU Architecture", runtime.GOARCH)
	serverTags.Set("Operating System", runtime.GOOS)
	serverTags.Set("CPU Count", fmt.Sprint(runtime.NumCPU()))

	appTags := NewOrderedMap[string](false, 10)
	appTags.Set("Name", AppName)
	exec, _ := os.Executable()
	wind, exec := StringSplitLast(exec, '/', true)
	if exec == "" {
		_, exec = StringSplitLast(wind, '\\', true)
	}
	appTags.Set("Executable", exec)
	appTags.Set("Version", version)
	appTags.Set("Go Version", runtime.Version())
	appTags.Set("Compiler", runtime.Compiler)

	runtimeTags := NewOrderedMap[string](false, 10)
	args := os.Args
	if len(args) > 0 {
		args = args[1:]
	}
	runtimeTags.Set("Running Since", TimeRelative(&started))
	runtimeTags.Set("Arguments", StringJoin(args, " "))
	runtimeTags.Set("Max Processes", fmt.Sprint(runtime.GOMAXPROCS(-1)))
	runtimeTags.Set("Go Routines", fmt.Sprint(runtime.NumGoroutine()))
	runtimeTags.Set("Cgo Calls", fmt.Sprint(runtime.NumCgoCall()))

	return &DebugInfo{ServerTags: serverTags, AppTags: appTags, RuntimeTags: runtimeTags, Mods: bi.Deps}
}

func DebugBuildInfo() *debug.BuildInfo {
	mods, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}
	return mods
}

func DebugStartCPUProfile() error {
	f, err := os.Create("./tmp/cpu.pprof")
	if err != nil {
		return err
	}
	return pprof.StartCPUProfile(f)
}

func DebugTakeHeapProfile() error {
	f, err := os.Create("./tmp/mem.pprof")
	if err != nil {
		return err
	}
	return pprof.WriteHeapProfile(f)
}

func DebugMemStats() *runtime.MemStats {
	x := &runtime.MemStats{}
	runtime.ReadMemStats(x)
	return x
}
