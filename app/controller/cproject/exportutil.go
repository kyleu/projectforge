package cproject

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const keyNew = "new"

func exportModelFromForm(frm util.ValueMap, m *model.Model) error {
	get := func(k string, def string) string {
		x := frm.GetStringOpt(k)
		if x == "" {
			return def
		}
		return x
	}
	m.Name = get("name", m.Name)
	m.Package = get("package", m.Package)
	m.Group = util.StringSplitAndTrim(get("group", strings.Join(m.Group, "/")), "/")
	m.Description = get("description", m.Description)
	m.Icon = get("icon", m.Icon)

	ords := filter.Orderings{}
	err := util.FromJSON([]byte(get("ordering", util.ToJSON(m.Ordering))), &ords)
	if err != nil {
		return errors.Wrap(err, "invalid ordering")
	}
	m.Ordering = ords

	sIdx, _ := strconv.ParseInt(get("sortIndex", "0"), 10, 64)
	m.SortIndex = int(sIdx)
	m.View = get("view", m.View)
	m.Search = util.StringSplitAndTrim(get("search", strings.Join(m.Search, ",")), ",")
	m.Tags = util.StringSplitAndTrim(get("tags", strings.Join(m.Tags, ",")), ",")
	m.TitleOverride = get("titleOverride", m.TitleOverride)
	m.PluralOverride = get("pluralOverride", m.PluralOverride)
	m.ProperOverride = get("properOverride", m.ProperOverride)
	m.TableOverride = get("tableOverride", m.TableOverride)
	m.RouteOverride = get("routeOverride", m.RouteOverride)

	cfg := util.ValueMap{}
	err = util.FromJSON([]byte(get("config", util.ToJSON(m.Config))), &cfg)
	if err != nil {
		return errors.Wrap(err, "invalid config")
	}
	m.Config = cfg

	cols := model.Columns{}
	err = util.FromJSON([]byte(get("columns", util.ToJSON(m.Columns))), &cols)
	if err != nil {
		return errors.Wrap(err, "invalid columns")
	}
	m.Columns = cols

	rels := model.Relations{}
	err = util.FromJSON([]byte(get("relations", util.ToJSON(m.Relations))), &rels)
	if err != nil {
		return errors.Wrap(err, "invalid relations")
	}
	m.Relations = rels

	idxs := model.Indexes{}
	err = util.FromJSON([]byte(get("indexes", util.ToJSON(m.Indexes))), &idxs)
	if err != nil {
		return errors.Wrap(err, "invalid indexes")
	}
	m.Indexes = idxs

	var sd [][]any
	err = util.FromJSON([]byte(get("seedData", util.ToJSON(m.SeedData))), &sd)
	if err != nil {
		return errors.Wrap(err, "invalid seed data")
	}
	m.SeedData = sd

	return nil
}

func exportLoadModel(r *http.Request, as *app.State, logger util.Logger) (*project.Project, *model.Model, *model.Args, error) {
	prj, args, err := exportLoadPrjArgs(r, as, logger)
	if err != nil {
		return nil, nil, nil, err
	}
	modelKey, err := cutil.RCRequiredString(r, "model", false)
	if err != nil {
		return nil, nil, nil, err
	}
	var mdl *model.Model
	if modelKey == keyNew {
		mdl = &model.Model{}
	} else {
		mdl = args.Models.Get(modelKey)
	}
	if mdl == nil {
		return nil, nil, nil, errors.Errorf("no model found with key [%s]", modelKey)
	}

	return prj, mdl, args, nil
}

func exportEnumFromForm(frm util.ValueMap, e *enum.Enum) error {
	get := func(k string, def string) string {
		x := frm.GetStringOpt(k)
		if x == "" {
			return def
		}
		return x
	}
	e.Name = get("name", e.Name)
	e.Package = get("package", e.Package)
	e.Group = util.StringSplitAndTrim(get("group", strings.Join(e.Group, "/")), "/")
	e.Description = get("description", e.Description)
	e.Icon = get("icon", e.Icon)

	e.Tags = util.StringSplitAndTrim(get("tags", strings.Join(e.Tags, ",")), ",")
	e.TitleOverride = get("titleOverride", e.TitleOverride)
	e.ProperOverride = get("properOverride", e.ProperOverride)

	valuesStr := get("values", util.ToJSON(e.Values))
	err := util.FromJSON([]byte(valuesStr), &e.Values)
	if err != nil {
		return err
	}

	cfg := util.ValueMap{}
	err = util.FromJSON([]byte(get("config", util.ToJSON(e.Config))), &cfg)
	if err != nil {
		return errors.Wrap(err, "invalid config")
	}
	e.Config = cfg
	return nil
}

func exportLoadEnum(r *http.Request, as *app.State, logger util.Logger) (*project.Project, *enum.Enum, *model.Args, error) {
	prj, args, err := exportLoadPrjArgs(r, as, logger)
	if err != nil {
		return nil, nil, nil, err
	}
	enumKey, err := cutil.RCRequiredString(r, "enum", false)
	if err != nil {
		return nil, nil, nil, err
	}
	var en *enum.Enum
	if enumKey == keyNew {
		en = &enum.Enum{}
	} else {
		en = args.Enums.Get(enumKey)
	}
	if en == nil {
		return nil, nil, nil, errors.Errorf("no model found with key [%s]", enumKey)
	}

	return prj, en, args, nil
}

func exportLoadPrjArgs(r *http.Request, as *app.State, logger util.Logger) (*project.Project, *model.Args, error) {
	prj, err := getProject(r, as)
	if err != nil {
		return nil, nil, err
	}
	args, err := prj.ModuleArgExport(as.Services.Projects, logger)
	if err != nil {
		return nil, nil, err
	}
	return prj, args, nil
}
