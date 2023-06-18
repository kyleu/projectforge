package action

import (
	"github.com/samber/lo"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/project/svg"
	"projectforge.dev/projectforge/app/util"
)

const keyTrue, keyTag, keyTags = "true", "tag", "tags"

func onRules(pm *PrjAndMods) *Result {
	ret := newResult(TypeRules, pm.Prj, pm.Cfg, pm.Logger)
	if pm.EArgs == nil {
		ret.Status = "OK"
		return ret
	}

	fs := pm.PSvc.GetFilesystem(pm.Prj)
	icons, err := svg.List(fs, pm.Logger)
	if err != nil {
		return ret.WithError(err)
	}
	forbidden := []string{"app", "check", "down", "edit", "left", "question", "right", "search", "searchbox", "times", "up", "star"}
	cleanIcons := lo.Filter(icons, func(x string, _ int) bool {
		return !slices.Contains(forbidden, x)
	})
	for _, m := range pm.EArgs.Models {
		for slices.Contains(forbidden, m.Icon) {
			idx := util.StringHash(m.Name) % uint32(len(cleanIcons))
			m.Icon = cleanIcons[idx]
		}
		m.AddTag("audit")
		m.AddTag("search")
		for _, col := range m.Columns {
			switch strings.ToLower(col.Name) {
			case "name", "title":
				if len(m.Columns.WithTag("title")) == 0 {
					col.AddTag("title")
				}
			}
			if col.PK {
				col.Search = true
				col.RemoveTag("search")
			}
		}
	}

	b, err := os.ReadFile("rules.json")
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
		m := pm.EArgs.Models.Get(split[0])
		if m == nil {
			return errors.Errorf("no model found with name [%s]", split[0])
		}
		switch split[1] {
		case "group":
			m.Group = util.StringSplitAndTrim(v, ".")
		case keyTag:
			m.AddTag(v)
		case keyTags:
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
			case keyTag:
				col.AddTag(v)
			case keyTags:
				for _, t := range util.StringSplitAndTrim(v, ",") {
					col.AddTag(t)
				}
			case "search":
				col.Search = v == keyTrue
			default:
				return errors.Errorf("unable to handle action [%s] for [%s.%s]", split[2], split[0], split[1])
			}
		}
	}
	return nil
}
