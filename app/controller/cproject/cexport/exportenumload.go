package cexport

import (
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cproject"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func exportEnumFromForm(frm util.ValueMap, e *enum.Enum) error {
	get := func(k string, def string) string {
		return util.OrDefault(frm.GetStringOpt(k), def)
	}
	e.Name = get("name", e.Name)
	e.Package = get("package", e.Package)
	e.Group = util.StringSplitAndTrim(get("group", util.StringJoin(e.Group, "/")), "/")
	e.Description = get("description", e.Description)
	e.Icon = get("icon", e.Icon)

	e.Tags = util.StringSplitAndTrim(get("tags", util.StringJoin(e.Tags, ",")), ",")
	e.TitleOverride = get("titleOverride", e.TitleOverride)
	e.ProperOverride = get("properOverride", e.ProperOverride)

	valuesStr := get("values", util.ToJSON(e.Values))
	err := util.FromJSON([]byte(valuesStr), &e.Values)
	if err != nil {
		return err
	}

	cfg, err := util.FromJSONMap([]byte(get("config", util.ToJSON(e.Config))))
	if err != nil {
		return errors.Wrap(err, "invalid config")
	}
	e.Config = cfg
	return nil
}

func exportLoadEnum(r *http.Request, as *app.State, logger util.Logger) (*project.Project, *enum.Enum, error) {
	prj, err := cproject.GetProjectWithArgs(r, as, logger)
	if err != nil {
		return nil, nil, err
	}
	enumKey, err := cutil.PathString(r, "enum", false)
	if err != nil {
		return nil, nil, err
	}
	var en *enum.Enum
	if enumKey == keyNew {
		en = &enum.Enum{}
	} else {
		en = prj.ExportArgs.Enums.Get(enumKey)
	}
	if en == nil {
		return nil, nil, errors.Errorf("no model found with key [%s]", enumKey)
	}

	return prj, en, nil
}
