package settings

import (
	"fmt"
	"runtime"
	"runtime/pprof"

	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/util"
)

func runGC(_ *mvc.State, _ *mvc.PageState) ([]string, error) {
	t := util.TimerStart()
	runtime.GC()
	return []string{"ran garbage collection", fmt.Sprintf("elapsed: %s", t.EndString())}, nil
}

func runHeapDump(_ *mvc.State, _ *mvc.PageState) ([]string, error) {
	if err := util.DebugTakeHeapProfile(); err != nil {
		return nil, err
	}
	return []string{"wrote heap profile", "path: ./tmp/mem.pprof"}, nil
}

func runCPUStart(_ *mvc.State, _ *mvc.PageState) ([]string, error) {
	if err := util.DebugStartCPUProfile(); err != nil {
		return nil, err
	}
	return []string{"started CPU profile", "path: ./tmp/cpu.pprof"}, nil
}

func runCPUStop(_ *mvc.State, _ *mvc.PageState) ([]string, error) {
	pprof.StopCPUProfile()
	return []string{"stopped CPU profile"}, nil
}
