package project

import (
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var (
	templateServicesOrder = []string{"audit", "user", "har", "process", "notebook", "scripting", "websocket", "help"}
	templateServicesNames = map[string]string{
		"audit": "Audit", "har": "Har", "process": "Exec", "scripting": "Script", "notebook": "Notebook",
		"websocket": "Socket", "user": "User", "help": "Help",
	}
	templateServicesKeys = map[string]string{
		"audit": "audit", "har": "har", "process": "exec", "scripting": "scripting", "notebook": "notebook",
		"websocket": "websocket", "user": "user", "help": "help",
	}
	templateServicesRefs = map[string]string{
		"audit":     "auditSvc",
		"har":       "har.NewService(st.Files)",
		"notebook":  "notebook.NewService()",
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
			if key == "user" && t.ExportArgs != nil && t.ExportArgs.Models.Get("user") != nil {
				continue
			}
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
	svcs, names, maxNameLength := t.servicesNames()
	if len(svcs) == 0 {
		return "type CoreServices struct{}"
	}
	ret := util.NewStringSlice([]string{"type CoreServices struct {"})
	types := lo.Map(svcs, func(svc string, _ int) string {
		return "*" + templateServicesKeys[svc] + ".Service"
	})
	for idx := range svcs {
		ret.Pushf("\t%s %s", util.StringPad(names[idx], maxNameLength), types[idx])
	}
	ret.Push("}")
	return ret.Join("\n")
}

func (t *TemplateContext) ServicesImports() string {
	svcs, _, _ := t.servicesNames()
	ret := &util.StringSlice{}
	for _, svc := range svcs {
		if svc == "user" {
			continue
		}
		ret.Pushf("\t\"%s/app/lib/%s\"", t.Package, templateServicesKeys[svc])
	}
	ret.Sort()
	if slices.Contains(svcs, "user") {
		ret.Pushf("\t\"%s/app/user\"", t.Package)
	}
	ret.Pushf("\t\"%s/app/util\"", t.Package)
	return strings.TrimPrefix(ret.Join("\n"), "\t")
}

func (t *TemplateContext) ServicesConstructor() string {
	svcs, names, maxNameLength := t.servicesNames()
	if len(svcs) == 0 {
		return "CoreServices{}"
	}
	ret := util.NewStringSlice([]string{"CoreServices{"})
	refs := lo.Map(svcs, func(svc string, _ int) string {
		return templateServicesRefs[svc]
	})
	args := ""
	if t.HasModule("audit") {
		args += ", auditSvc"
	}
	for idx, key := range svcs {
		if key == "user" && t.ExportArgs != nil && t.ExportArgs.Models.Get("user") != nil {
			continue
		}
		ret.Pushf("\t\t%s %s,", util.StringPad(names[idx]+":", maxNameLength+1), refs[idx])
	}
	ret.Push("\t}")
	return ret.Join("\n")
}
