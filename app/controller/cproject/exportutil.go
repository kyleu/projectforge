package cproject

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

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

	m.Search = util.StringSplitAndTrim(get("search", strings.Join(m.Search, ",")), ",")
	m.History = get("history", m.History)
	m.Tags = util.StringSplitAndTrim(get("tags", strings.Join(m.Tags, ",")), ",")
	m.TitleOverride = get("titleOverride", m.TitleOverride)
	m.ProperOverride = get("propoerOverride", m.ProperOverride)

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
		return errors.Wrap(err, "invalid indexs")
	}
	m.Indexes = idxs

	return nil
}

func exportLoad(rc *fasthttp.RequestCtx, as *app.State, logger util.Logger) (*project.Project, *model.Model, *model.Args, error) {
	prj, err := getProject(rc, as)
	if err != nil {
		return nil, nil, nil, err
	}

	modelKey, err := cutil.RCRequiredString(rc, "model", false)
	if err != nil {
		return nil, nil, nil, err
	}

	args, err := prj.ModuleArgExport(as.Services.Projects, logger)
	if err != nil {
		return nil, nil, nil, err
	}
	var mdl *model.Model
	if modelKey == "new" {
		mdl = &model.Model{}
	} else {
		mdl = args.Models.Get(modelKey)
	}

	return prj, mdl, args, nil
}
