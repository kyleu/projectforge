package cexport

import (
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cproject"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func exportEventFromForm(frm util.ValueMap, m *model.Event) error {
	var err error
	get := func(k string, def string) string {
		return util.OrDefault(frm.GetStringOpt(k), def)
	}
	m.Name = get("name", m.Name)
	m.Package = get("package", m.Package)
	m.Group = util.StringSplitAndTrim(get("group", util.StringJoin(m.Group, "/")), "/")
	m.Schema = get("schema", m.Schema)
	m.Description = get("description", m.Description)
	m.Icon = get("icon", m.Icon)

	m.Tags = util.StringSplitAndTrim(get("tags", util.StringJoin(m.Tags, ",")), ",")
	m.TitleOverride = get("titleOverride", m.TitleOverride)
	m.PluralOverride = get("pluralOverride", m.PluralOverride)
	m.ProperOverride = get("properOverride", m.ProperOverride)

	m.Config, err = util.FromJSONMap([]byte(get("config", util.ToJSON(m.Config))))
	if err != nil {
		return errors.Wrap(err, "invalid config")
	}

	m.Columns, err = util.FromJSONObj[model.Columns]([]byte(get("columns", util.ToJSON(m.Columns))))
	if err != nil {
		return errors.Wrap(err, "invalid columns")
	}

	m.Imports, err = util.FromJSONObj[model.Imports]([]byte(get("imps", util.ToJSON(m.Imports))))
	if err != nil {
		return errors.Wrap(err, "invalid imports")
	}

	return nil
}

func exportLoadEvent(r *http.Request, as *app.State, logger util.Logger) (*project.Project, *model.Event, error) {
	prj, err := cproject.GetProjectWithArgs(r, as, logger)
	if err != nil {
		return nil, nil, err
	}
	eventKey, err := cutil.PathString(r, "event", false)
	if err != nil {
		return nil, nil, err
	}
	var mdl *model.Event
	if eventKey == keyNew {
		mdl = &model.Event{}
	} else {
		mdl = prj.ExportArgs.Events.Get(eventKey)
	}
	if mdl == nil {
		return nil, nil, errors.Errorf("no event found with key [%s]", eventKey)
	}

	return prj, mdl, nil
}
