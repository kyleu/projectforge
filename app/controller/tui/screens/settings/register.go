package settings

import (
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/lib/menu"
)

func Register(reg *screens.Registry) {
	reg.RegisterBefore(
		screens.KeyAbout,
		&menu.Item{Key: screens.KeySettings, Title: "Settings", Description: "Runtime and diagnostics", Icon: "settings", Route: screens.KeySettings},
		newMenuScreen(),
	)
	registerCoreScreens(reg)
	registerExecScreens(reg)
	registerTaskScreens(reg)
	registerSocketScreens(reg)
}

func registerCoreScreens(reg *screens.Registry) {
	reg.AddScreen(newDataScreen(keyServer, "App Information", serverLines))
	reg.AddScreen(newModulesScreen())
	reg.AddScreen(newDataScreen(keyMemUsage, "Memory Usage", memUsageLines))
	reg.AddScreen(newRoutesScreen())
	reg.AddScreen(newDataScreen(keySitemap, "Sitemap", sitemapLines))
	reg.AddScreen(newActionScreen(keyGC, "Collect Garbage", runGC))
	reg.AddScreen(newActionScreen(keyHeap, "Write Memory Dump", runHeapDump))
	reg.AddScreen(newActionScreen(keyCPUStart, "Start CPU Profile", runCPUStart))
	reg.AddScreen(newActionScreen(keyCPUStop, "Stop CPU Profile", runCPUStop))
}
