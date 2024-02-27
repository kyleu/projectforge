// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"
)

var DEBUG = false

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
	serverTags.Append("Machine Name", hostname)
	serverTags.Append("CPU Architecture", runtime.GOARCH)
	serverTags.Append("Operating System", runtime.GOOS)
	serverTags.Append("CPU Count", fmt.Sprint(runtime.NumCPU()))

	appTags := NewOrderedMap[string](false, 10)
	appTags.Append("Name", AppName)
	exec, _ := os.Executable()
	wind, exec := StringSplitLast(exec, '/', true)
	if exec == "" {
		_, exec = StringSplitLast(wind, '\\', true)
	}
	appTags.Append("Executable", exec)
	appTags.Append("Version", version)
	appTags.Append("Go Version", runtime.Version())
	appTags.Append("Compiler", runtime.Compiler)

	runtimeTags := NewOrderedMap[string](false, 10)
	args := os.Args
	if len(args) > 0 {
		args = args[1:]
	}
	runtimeTags.Append("Running Since", TimeRelative(&started))
	runtimeTags.Append("Arguments", strings.Join(args, " "))
	runtimeTags.Append("Max Processes", fmt.Sprint(runtime.GOMAXPROCS(-1)))
	runtimeTags.Append("Go Routines", fmt.Sprint(runtime.NumGoroutine()))
	runtimeTags.Append("Cgo Calls", fmt.Sprint(runtime.NumCgoCall()))

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
