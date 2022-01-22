package model

import (
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

type Args struct {
	Config  util.ValueMap `json:"config,omitempty"`
	Models  Models        `json:"models,omitempty"`
	Modules []string      `json:"-"`
}

func (a *Args) HasModule(key string) bool {
	return util.StringArrayContains(a.Modules, key)
}

func (a *Args) Validate() error {
	for _, m := range a.Models {
		err := m.Validate()
		if err != nil {
			return errors.Wrap(err, "invalid model ["+m.Name+"]")
		}
		for _, rel := range m.Relations {
			relTable := a.Models.Get(rel.Table)
			if relTable == nil {
				return errors.Errorf("relation [%s] refers to missing table [%s]", rel.Name, rel.Table)
			}
		}
	}
	return nil
}
