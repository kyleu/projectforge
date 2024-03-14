package project

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var (
	templateServicesOrder = []string{"audit", "user", "har", "process", "scripting", "websocket", "help"}
	templateServicesNames = map[string]string{
		"audit": "Audit", "har": "Har", "process": "Exec", "scripting": "Script", "websocket": "Socket", "user": "User", "help": "Help",
	}
	templateServicesKeys = map[string]string{
		"audit": "audit", "har": "har", "process": "exec", "scripting": "scripting", "websocket": "websocket", "user": "user", "help": "help",
	}
	templateServicesRefs = map[string]string{
		"audit":     "auditSvc",
		"har":       "har.NewServices(st.Files)",
		"process":   "exec.NewService()",
		"scripting": "scripting.NewService(st.Files, \"scripts\")",
		"user":      "user.NewService(st.Files, logger)",
		"websocket": "websocket.NewService(nil, nil, nil)",
		"help":      "help.NewService(logger)",
	}
)

func (t *TemplateContext) servicesNames() ([]string, []string, int) {
	var svcs []string
	for _, key := range templateServicesOrder {
		if t.HasModule(key) {
			svcs = append(svcs, key)
		}
	}
	names := lo.Map(svcs, func(svc string, _ int) string {
		return templateServicesNames[svc]
	})
	maxNameLength := util.StringArrayMaxLength(names)
	return svcs, names, maxNameLength
}

func (t *TemplateContext) ServicesDefinition() string {
	ret := []string{"type Services struct {"}
	w := func(msg string, args ...any) {
		ret = append(ret, fmt.Sprintf(msg, args...))
	}

	svcs, names, maxNameLength := t.servicesNames()
	types := lo.Map(svcs, func(svc string, _ int) string {
		return "*" + templateServicesKeys[svc] + ".Service"
	})

	if t.HasModule("export") {
		w("\tGeneratedServices")
		w("")
	}
	w("\t// add your dependencies here")
	w("")
	for idx := range svcs {
		w("\t%s %s", util.StringPad(names[idx], maxNameLength), types[idx])
	}
	w("}")
	return strings.Join(ret, "\n")
}

func (t *TemplateContext) ServicesImports() string {
	ret := []string{}
	w := func(msg string, args ...any) {
		ret = append(ret, fmt.Sprintf(msg, args...))
	}

	if t.HasModule("migration") {
		w("\t\"github.com/pkg/errors\"")
		w("")
	}

	var keys []string
	for _, key := range templateServicesOrder {
		if t.HasModule(key) {
			keys = append(keys, templateServicesKeys[key])
		}
	}
	if t.HasModule("migration") {
		keys = append(keys, "database/migrate")
	}
	slices.Sort(keys)
	for _, svc := range keys {
		w("\t\"%s/app/lib/%s\"", t.Package, svc)
	}
	w("\t\"%s/app/util\"", t.Package)
	if t.HasModule("migration") {
		w("\t\"%s/queries/migrations\"", t.Package)
	}
	return strings.TrimPrefix(strings.Join(ret, "\n"), "\t")
}

func (t *TemplateContext) ServicesConstructor() string {
	ret := []string{"&Services{"}
	w := func(msg string, args ...any) {
		ret = append(ret, fmt.Sprintf(msg, args...))
	}

	svcs, names, maxNameLength := t.servicesNames()
	refs := lo.Map(svcs, func(svc string, _ int) string {
		return templateServicesRefs[svc]
	})

	if t.HasModule("export") {
		args := ""
		if t.HasModule("readonlydb") {
			args += ", st.DBRead"
		}
		if t.HasModule("audit") {
			args += ", auditSvc"
		}
		w("\t\tGeneratedServices: initGeneratedServices(ctx, st.DB%s, logger),", args)
		w("")
	}
	for idx := range svcs {
		w("\t\t%s %s,", util.StringPad(names[idx]+":", maxNameLength+1), refs[idx])
	}
	w("\t}")
	return strings.Join(ret, "\n")
}
