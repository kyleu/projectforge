package search

import (
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
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
	return lo.Reject(p.Parts(), func(x string, _ int) bool {
		return strings.Contains(x, ":")
	})
}

func (p *Params) Keyed() util.ValueMap {
	x := lo.Filter(p.Parts(), func(x string, _ int) bool {
		return strings.Contains(x, ":")
	})
	ret := make(util.ValueMap, len(x))
	lo.ForEach(x, func(s string, _ int) {
		k, v := util.StringSplit(s, ':', true)
		ret[k] = v
	})
	return ret
}
