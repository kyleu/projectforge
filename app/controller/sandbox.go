package controller

import (
	"os"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
)

func Sandbox(rc *fasthttp.RequestCtx) {
	Act("sandbox", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", true)
		if err != nil {
			return "", err
		}
		ps.Data = "OK"
		err = cleanProject(as, key, ps.Logger)
		if err != nil {
			return "", err
		}
		return Render(rc, as, &views.Debug{}, ps)
	})
}

func cleanProject(st *app.State, prjKey string, logger util.Logger) error {
	svcs := st.Services
	prj, err := svcs.Projects.Get(prjKey)
	if err != nil {
		return err
	}
	args, err := prj.ModuleArgExport(svcs.Projects, logger)
	if err != nil {
		return err
	}

	for _, m := range args.Models {
		m.AddTag("search")
		for _, col := range m.Columns {
			if col.PK {
				col.Search = true
				col.RemoveTag("search")
			}
		}
	}

	b, err := os.ReadFile("rules.json")
	if err != nil {
		return err
	}
	rules := map[string]string{}
	err = util.FromJSON(b, &rules)
	if err != nil {
		return err
	}
	for k, v := range rules {
		split := util.StringSplitAndTrim(k, ".")
		if split[0] == "disabled" {
			continue
		}
		m := args.Models.Get(split[0])
		if m == nil {
			return errors.Errorf("no model found with name [%s]", split[0])
		}
		switch split[1] {
		case "group":
			m.Group = util.StringSplitAndTrim(v, ".")
		case "tag":
			m.AddTag(v)
		case "tags":
			for _, t := range util.StringSplitAndTrim(v, ",") {
				m.AddTag(t)
			}
		default:
			col := m.Columns.Get(split[1])
			if col == nil {
				return errors.Errorf("no column found with name [%s] in model [%s]", split[1], split[0])
			}
			switch split[2] {
			case "display":
				col.Display = v
			case "format":
				col.Format = v
			case "tag":
				col.AddTag(v)
			case "tags":
				for _, t := range util.StringSplitAndTrim(v, ",") {
					col.AddTag(t)
				}
			case "search":
				col.Search = v == "true"
			default:
				return errors.Errorf("unable to handle action [%s] for [%s.%s]", split[2], split[0], split[1])
			}
		}
	}

	err = svcs.Projects.Save(prj, logger)
	if err != nil {
		return err
	}

	return nil
}
