package action

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

const keyTag, keyTags = "tag", "tags"

func onRules(pm *PrjAndMods) *Result {
	ret := newResult(TypeRules, pm.Prj, pm.Cfg, pm.Logger)
	if pm.Prj.ExportArgs == nil {
		ret.Status = util.OK
		return ret
	}

	lo.ForEach(pm.Prj.ExportArgs.Models, func(m *model.Model, _ int) {
		m.AddTag("search")
		lo.ForEach(m.Columns, func(col *model.Column, _ int) {
			switch strings.ToLower(col.Name) {
			case "name", "title":
				if len(m.Columns.WithTag("title")) == 0 {
					col.AddTag("title")
				}
			}
			if col.PK {
				col.Search = true
			}
		})
	})

	fs, _ := filesystem.NewFileSystem(".", false, "")
	b, err := fs.ReadFile("rules.json")
	if err != nil {
		return ret.WithError(err)
	}
	rules := map[string]string{}
	err = util.FromJSON(b, &rules)
	if err != nil {
		return ret.WithError(err)
	}

	err = applyRules(pm, rules)
	if err != nil {
		return ret.WithError(err)
	}

	err = pm.PSvc.Save(pm.Prj, pm.Logger)
	if err != nil {
		return ret.WithError(err)
	}

	return ret
}

func applyRules(pm *PrjAndMods, rules map[string]string) error {
	for k, v := range rules {
		split := util.StringSplitAndTrim(k, ".")
		if split[0] == "disabled" {
			continue
		}
		m := pm.Prj.ExportArgs.Models.Get(split[0])
		if m == nil {
			return errors.Errorf("no model found with name [%s]", split[0])
		}
		switch split[1] {
		case "group":
			m.Group = util.StringSplitAndTrim(v, ".")
		case "schema":
			m.Schema = v
		case keyTag:
			m.AddTag(v)
		case keyTags:
			lo.ForEach(util.StringSplitAndTrim(v, ","), func(t string, _ int) {
				m.AddTag(t)
			})
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
			case "example":
				col.Example = v
			case "json":
				col.JSON = v
			case "validation":
				col.Validation = v
			case keyTag:
				col.AddTag(v)
			case keyTags:
				lo.ForEach(util.StringSplitAndTrim(v, ","), func(t string, _ int) {
					col.AddTag(t)
				})
			case "search":
				col.Search = v == util.BoolTrue
			default:
				return errors.Errorf("unable to handle action [%s] for [%s.%s]", split[2], split[0], split[1])
			}
		}
	}
	return nil
}
