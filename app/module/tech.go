package module

import "projectforge.dev/projectforge/app/util"

var techLookup = map[string][]string{
	"afero":         {"afero", "https://github.com/spf13/afero"},
	"android":       {"Android", "https://foo.com"},
	"cel":           {"CEL", "https://foo.com"},
	"coral":         {"Coral", "https://foo.com"},
	"gamut":         {"gamut", "https://github.com/muesli/gamut"},
	"go":            {"Golang", "https://foo.com"},
	"gocron":        {"gocron", "https://github.com/go-co-op/gocron/v2"},
	"goja":          {"goja", "https://github.com/dop251/goja"},
	"golang-cross":  {"golang-cross", "https://foo.com"},
	"gorilla":       {"gorilla", "https://foo.com"},
	"goth":          {"Goth", "https://github.com/markbates/goth"},
	"graphql-go":    {"graphql-go", "https://foo.com"},
	"ios":           {"iOS", "https://foo.com"},
	"markdown":      {"Markdown", "https://foo.com"},
	"mysql":         {"MySQL", "https://www.mysql.org"},
	"notarytool":    {"notarytool", "https://foo.com"},
	"observable":    {"Observable", "https://foo.com"},
	"opentelemetry": {"OpenTelemetry", "https://foo.com"},
	"playwright":    {"Playwright", "https://foo.com"},
	"postgres":      {"PostgreSQL", "https://www.postgresql.org"},
	"projectforge":  {"Project Forge", "https://projectforge.dev", "app"},
	"simple-icons":  {"Simple Icons", "https://foo.com"},
	"sqlite":        {"SQLite", "https://www.sqlite.org"},
	"sqlserver":     {"SQL Server", "https://www.microsoft.com/en-us/sql-server"},
	"swagger":       {"Swagger", "https://foo.com"},
	"typescript":    {"TypeScript", "https://foo.com"},
	"wasm":          {"WebAssembly", "https://foo.com"},
	"websocket":     {"WebSocket", "https://foo.com"},
	"xcode":         {"Xcode", "https://foo.com"},
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
