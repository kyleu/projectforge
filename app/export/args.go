package export

import (
	"github.com/kyleu/projectforge/app/util"
)

type Args struct {
	Models  Models   `json:"models"`
	Modules []string `json:"-"`
}

func (a *Args) HasModule(key string) bool {
	return util.StringArrayContains(a.Modules, key)
}
