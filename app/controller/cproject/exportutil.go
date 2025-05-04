package cproject

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const keyNew = "new"

func exportModelFromForm(frm util.ValueMap, m *model.Model) error {
	get := func(k string, def string) string {
		return util.OrDefault(frm.GetStringOpt(k), def)
	}
	m.Name = get("name", m.Name)
	m.Package = get("package", m.Package)
	m.Group = util.StringSplitAndTrim(get("group", util.StringJoin(m.Group, "/")), "/")
	m.Schema = get("schema", m.Schema)
	m.Description = get("description", m.Description)
	m.Icon = get("icon", m.Icon)

	ords, err := util.FromJSONObj[filter.Orderings]([]byte(get("ordering", util.ToJSON(m.Ordering))))
	if err != nil {
		return errors.Wrap(err, "invalid ordering")
	}
	m.Ordering = ords

	sIdx, _ := strconv.ParseInt(get("sortIndex", "0"), 10, 64)
	m.SortIndex = int(sIdx)
	m.View = get("view", m.View)
	m.Search = util.StringSplitAndTrim(get("search", util.StringJoin(m.Search, ",")), ",")
	m.Tags = util.StringSplitAndTrim(get("tags", util.StringJoin(m.Tags, ",")), ",")
	m.TitleOverride = get("titleOverride", m.TitleOverride)
	m.PluralOverride = get("pluralOverride", m.PluralOverride)
	m.ProperOverride = get("properOverride", m.ProperOverride)
	m.TableOverride = get("tableOverride", m.TableOverride)
	m.RouteOverride = get("routeOverride", m.RouteOverride)

	m.Config, err = util.FromJSONMap([]byte(get("config", util.ToJSON(m.Config))))
	if err != nil {
		return errors.Wrap(err, "invalid config")
	}

	m.Columns, err = util.FromJSONObj[model.Columns]([]byte(get("columns", util.ToJSON(m.Columns))))
	if err != nil {
		return errors.Wrap(err, "invalid columns")
	}

	m.Relations, err = util.FromJSONObj[model.Relations]([]byte(get("relations", util.ToJSON(m.Relations))))
	if err != nil {
		return errors.Wrap(err, "invalid relations")
	}

	m.Indexes, err = util.FromJSONObj[model.Indexes]([]byte(get("indexes", util.ToJSON(m.Indexes))))
	if err != nil {
		return errors.Wrap(err, "invalid indexes")
	}

	m.SeedData, err = util.FromJSONObj[[][]any]([]byte(get("seedData", util.ToJSON(m.SeedData))))
	if err != nil {
		return errors.Wrap(err, "invalid seed data")
	}

	m.Links, err = util.FromJSONObj[model.Links]([]byte(get("links", util.ToJSON(m.Links))))
	if err != nil {
		return errors.Wrap(err, "invalid links")
	}

	m.Imports, err = util.FromJSONObj[model.Imports]([]byte(get("imps", util.ToJSON(m.Imports))))
	if err != nil {
		return errors.Wrap(err, "invalid imports")
	}

	return nil
}

func exportLoadModel(r *http.Request, as *app.State, logger util.Logger) (*project.Project, *model.Model, error) {
	prj, err := getProjectWithArgs(r, as, logger)
	if err != nil {
		return nil, nil, err
	}
	modelKey, err := cutil.PathString(r, "model", false)
	if err != nil {
		return nil, nil, err
	}
	var mdl *model.Model
	if modelKey == keyNew {
		mdl = &model.Model{}
	} else {
		mdl = prj.ExportArgs.Models.Get(modelKey)
	}
	if mdl == nil {
		return nil, nil, errors.Errorf("no model found with key [%s]", modelKey)
	}

	return prj, mdl, nil
}

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
	prj, err := getProjectWithArgs(r, as, logger)
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
