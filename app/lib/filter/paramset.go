package filter

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type ParamSet map[string]*Params

func (s ParamSet) Get(key string, allowed []string, logger util.Logger) *Params {
	x, ok := s[key]
	if !ok {
		return &Params{Key: key}
	}
	return x.Filtered(key, allowed, logger).Sanitize(key)
}

func (s ParamSet) Sanitized(key string, logger util.Logger, defaultOrderings ...*Ordering) *Params {
	ret := s.Get(key, nil, logger).Sanitize(key, defaultOrderings...)
	return ret
}

func (s ParamSet) Specifies(key string) bool {
	x, ok := s[key]
	if !ok {
		return false
	}
	return !x.IsDefault()
}

func (s ParamSet) String() string {
	return util.StringJoin(lo.Map(lo.Values(s), func(p *Params, _ int) string {
		return p.String()
	}), ", ")
}
