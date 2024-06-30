package module

import "projectforge.dev/projectforge/app/util"

var techLookup = map[string][]string{
	"afero":         {"afero", "https://github.com/spf13/afero", "folder"},
	"android":       {"Android", "https://android.com"},
	"cel":           {"CEL", "https://github.com/google/cel-spec"},
	"coral":         {"Coral", "https://github.com/muesli/coral", "terminal"},
	"gamut":         {"gamut", "https://github.com/muesli/gamut", "paintbrush"},
	"go":            {"Golang", "https://go.dev"},
	"gocron":        {"gocron", "https://github.com/go-co-op/gocron/v2", "hourglass"},
	"goja":          {"goja", "https://github.com/dop251/goja", "code"},
	"golang-cross":  {"golang-cross", "https://github.com/gythialy/golang-cross", "dna"},
	"gorilla":       {"gorilla", "https://gorilla.github.io", "shield"},
	"goth":          {"Goth", "https://github.com/markbates/goth", "users"},
	"graphql-go":    {"graphql-go", "https://github.com/graphql-go/graphql", "graphql"},
	"ios":           {"iOS", "https://www.apple.com/ios", "apple"},
	"markdown":      {"Markdown", "https://en.wikipedia.org/wiki/Markdown"},
	"mysql":         {"MySQL", "https://www.mysql.org"},
	"notarytool":    {"notarytool", "https://developer.apple.com", "apple"},
	"observable":    {"Observable", "https://observablehq.com/framework", "book"},
	"opentelemetry": {"OpenTelemetry", "https://opentelemetry.io"},
	"playwright":    {"Playwright", "https://playwright.dev"},
	"postgres":      {"PostgreSQL", "https://www.postgresql.org"},
	"projectforge":  {"Project Forge", "https://projectforge.dev", "app"},
	"simple-icons":  {"Simple Icons", "https://simpleicons.org", "simpleicons"},
	"sqlite":        {"SQLite", "https://www.sqlite.org"},
	"sqlserver":     {"SQL Server", "https://www.microsoft.com/en-us/sql-server"},
	"swagger":       {"Swagger", "https://swagger.io", "print"},
	"typescript":    {"TypeScript", "https://www.typescriptlang.org"},
	"wasm":          {"WebAssembly", "https://webassembly.org"},
	"websocket":     {"WebSocket", "https://developer.mozilla.org", "satellite"},
	"xcode":         {"Xcode", "https://developer.apple.com/xcode", "apple"},
}

func Tech(k string, logger util.Logger) (string, string, string) {
	x, ok := techLookup[k]
	if !ok {
		logger.Warnf("missing technology key [%s]", k)
		return k, "UNKNOWN", "UNKNOWN"
	}
	if len(x) > 2 {
		return x[0], x[1], x[2]
	}
	icon := "question"
	if _, ok := util.SVGLibrary[k]; ok {
		icon = k
	} else {
		logger.Warnf("missing icon [%s]", k)
	}
	return x[0], x[1], icon
}
