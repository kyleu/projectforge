package template

import (
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var (
	templateServicesOrder = []string{"audit", keyUser, "har", "process", "notebook", "schedule", "scripting", "websocket", "task", "help"}
	templateServicesNames = map[string]string{
		"audit": "Audit", "har": "Har", "notebook": "Notebook", "process": "Exec", "schedule": "Schedule",
		"scripting": "Script", "task": "Task", "websocket": "Socket", keyUser: "User", "help": "Help",
	}
	templateServicesKeys = map[string]string{
		"audit": "audit", "har": "har", "notebook": "notebook", "process": "exec", "schedule": "schedule",
		"scripting": "scripting", "task": "task", "websocket": "websocket", keyUser: keyUser, "help": "help",
	}
	templateServicesRefs = map[string]string{
		"audit":     "auditSvc",
		"har":       "har.NewService(st.Files)",
		"notebook":  "notebook.NewService()",
		"process":   "exec.NewService()",
		"schedule":  "schedule.NewService()",
		"scripting": "scripting.NewService(st.Files, \"scripts\")",
		"task":      "task.NewService(st.Files, \"task_history\")",
		keyUser:     "user.NewService(st.Files, logger)",
		"websocket": "websocket.NewService(nil, nil)",
		"help":      "help.NewService(logger)",
	}
)

func (t *TemplateContext) servicesNames() ([]string, []string, int) {
	var svcs []string
	for _, key := range templateServicesOrder {
		if t.HasModule(key) {
			if key == keyUser && t.ExportArgs != nil && t.ExportArgs.Models.Get(keyUser) != nil {
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
		if svc == keyUser {
			continue
		}
		ret.Pushf("\t\"%s/app/lib/%s\"", t.Package, templateServicesKeys[svc])
	}
	ret.Sort()
	if slices.Contains(svcs, keyUser) {
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
	for idx, key := range svcs {
		if key == keyUser && t.ExportArgs != nil && t.ExportArgs.Models.Get(keyUser) != nil {
			continue
		}
		ret.Pushf("\t\t%s %s,", util.StringPad(names[idx]+":", maxNameLength+1), refs[idx])
	}
	ret.Push("\t}")
	return ret.Join("\n")
}
