package settings

type optionKind string

const (
	kindScreen optionKind = "screen"
	kindAction optionKind = "action"
)

type adminOption struct {
	Key         string
	Title       string
	Description string
	Route       string
	Kind        optionKind
}

const (
	keyServer   = "settings.admin.server"
	keyModules  = "settings.admin.modules"
	keyExec     = "settings.admin.exec"
	keyTask     = "settings.admin.task"
	keySitemap  = "settings.admin.sitemap"
	keyRoutes   = "settings.admin.routes"
	keySockets  = "settings.admin.sockets"
	keyMemUsage = "settings.admin.memusage"
	keyGC       = "settings.admin.gc"
	keyHeap     = "settings.admin.heap"
	keyCPUStart = "settings.admin.cpu.start"
	keyCPUStop  = "settings.admin.cpu.stop"
)

var adminOptions = []adminOption{
	{Key: keyServer, Title: "App Information", Description: "All sorts of info about the server and runtime", Route: "/admin/server", Kind: kindScreen},
	{Key: keyModules, Title: "Go Modules", Description: "The Go modules used by Project Forge", Route: "/admin/modules", Kind: kindScreen},
	{Key: keyExec, Title: "Managed Processes", Description: "Manage OS processes from within this app", Route: "/admin/exec", Kind: kindScreen},
	{Key: keyTask, Title: "Task Engine", Description: "See the tasks that have run, and start new runs", Route: "/admin/task", Kind: kindScreen},
	{Key: keySitemap, Title: "Sitemap", Description: "Displays available HTTP actions", Route: "/admin/sitemap", Kind: kindScreen},
	{Key: keyRoutes, Title: "HTTP Routes", Description: "Enumerates all registered HTTP routes", Route: "/admin/routes", Kind: kindScreen},
	{Key: keySockets, Title: "Active WebSockets", Description: "Manage the active WebSockets in this server", Route: "/admin/sockets", Kind: kindScreen},
	{Key: keyMemUsage, Title: "Memory Usage", Description: "Detailed memory usage statistics", Route: "/admin/memusage", Kind: kindScreen},
	{Key: keyGC, Title: "Collect Garbage", Description: "Runs the Go garbage collector", Route: "/admin/gc", Kind: kindAction},
	{Key: keyHeap, Title: "Write Memory Dump", Description: "Writes a memory dump to ./tmp/mem.pprof", Route: "/admin/heap", Kind: kindAction},
	{Key: keyCPUStart, Title: "Start CPU Profile", Description: "Profiles CPU using ./tmp/cpu.pprof", Route: "/admin/cpu/start", Kind: kindAction},
	{Key: keyCPUStop, Title: "Stop CPU Profile", Description: "Stops the active CPU profile", Route: "/admin/cpu/stop", Kind: kindAction},
}
