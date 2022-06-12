package model

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/util"
)

type Args struct {
	Config  util.ValueMap `json:"config,omitempty"`
	Models  Models        `json:"models,omitempty"`
	Groups  Groups        `json:"groups,omitempty"`
	Modules []string      `json:"-"`
}

func (a *Args) HasModule(key string) bool {
	return slices.Contains(a.Modules, key)
}

func (a *Args) DBRef() string {
	if a.HasModule("readonlydb") {
		return "dbRead"
	}
	return "db"
}

func (a *Args) Validate() error {
	packages := make(map[string]struct{}, len(a.Models))
	for _, m := range a.Models {
		err := m.Validate(a.Modules)
		if err != nil {
			return errors.Wrap(err, "invalid model ["+m.Name+"]")
		}
		for _, rel := range m.Relations {
			relTable := a.Models.Get(rel.Table)
			if relTable == nil {
				return errors.Errorf("relation [%s] refers to missing table [%s]", rel.Name, rel.Table)
			}
			for _, t := range rel.Tgt {
				if relTable.Columns.Get(t) == nil {
					return errors.Errorf("relation [%s] references missing target column [%s]", rel.Name, t)
				}
			}
		}
		if _, ok := packages[m.Package]; ok {
			return errors.Wrap(err, "multiple models are in package ["+m.Package+"]")
		}
		packages[m.Package] = struct{}{}
	}
	return nil
}
