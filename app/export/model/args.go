package model

import (
	"github.com/kyleu/projectforge/app/util"
)

type Args struct {
	Config  util.ValueMap `json:"config,omitempty"`
	Models  Models        `json:"models,omitempty"`
	Modules []string      `json:"-"`
}

func (a *Args) HasModule(key string) bool {
	return util.StringArrayContains(a.Modules, key)
}
