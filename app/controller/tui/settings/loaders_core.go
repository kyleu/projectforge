package settings

import (
	"fmt"
	"runtime"

	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/util"
)

func serverLines(ts *mvc.State, _ *mvc.PageState) ([]string, error) {
	info := util.DebugGetInfo(ts.App.BuildInfo.Version, ts.App.Started)
	ret := []string{"Server Information"}
	ret = appendTags(ret, info.ServerTags)
	ret = append(ret, "", "Runtime Information")
	ret = appendTags(ret, info.RuntimeTags)
	ret = append(ret, "", "App Information")
	ret = appendTags(ret, info.AppTags)
	return ret, nil
}

func memUsageLines(_ *mvc.State, _ *mvc.PageState) ([]string, error) {
	m := util.DebugMemStats()
	return []string{
		fmt.Sprintf("total allocations: %s", util.ByteSizeSI(int64(m.TotalAlloc))),
		fmt.Sprintf("system memory: %s", util.ByteSizeSI(int64(m.Sys))),
		fmt.Sprintf("heap allocated: %s", util.ByteSizeSI(int64(m.HeapAlloc))),
		fmt.Sprintf("heap in use: %s", util.ByteSizeSI(int64(m.HeapInuse))),
		fmt.Sprintf("stack in use: %s", util.ByteSizeSI(int64(m.StackInuse))),
		fmt.Sprintf("gc count: %d", m.NumGC),
		fmt.Sprintf("gc cpu pct: %.2f", m.GCCPUFraction*100),
		fmt.Sprintf("next gc target: %s", util.ByteSizeSI(int64(m.NextGC))),
		fmt.Sprintf("goroutines: %d", runtime.NumGoroutine()),
	}, nil
}

func appendTags(ret []string, tags *util.OrderedMap[string]) []string {
	for _, k := range tags.Order {
		ret = append(ret, fmt.Sprintf("%s: %s", k, tags.GetSimple(k)))
	}
	return ret
}
