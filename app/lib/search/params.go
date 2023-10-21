// Package search - Content managed by Project Forge, see [projectforge.md] for details.
package search

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/util"
)

type Params struct {
	Q  string          `json:"q"`
	PS filter.ParamSet `json:"ps,omitempty"`
}

func (p *Params) String() string {
	return p.Q
}

func (p *Params) Parts() []string {
	return util.StringSplitAndTrim(p.Q, " ")
}

func (p *Params) General() []string {
	return lo.Filter(p.Parts(), func(x string, _ int) bool {
		return !strings.Contains(x, ":")
	})
}

func (p *Params) Keyed() map[string]string {
	x := lo.Filter(p.Parts(), func(x string, _ int) bool {
		return strings.Contains(x, ":")
	})
	ret := make(map[string]string, len(x))
	lo.ForEach(x, func(s string, _ int) {
		k, v := util.StringSplit(s, ':', true)
		ret[k] = v
	})
	return ret
}
