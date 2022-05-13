// Content managed by Project Forge, see [projectforge.md] for details.
package filter

import (
	"strings"

	"projectforge.dev/projectforge/app/util"
)

type ParamSet map[string]*Params

func (s ParamSet) Get(key string, allowed []string, logger util.Logger) *Params {
	x, ok := s[key]
	if !ok {
		return &Params{Key: key}
	}
	return x.Filtered(allowed, logger).Sanitize(key)
}

func (s ParamSet) String() string {
	ret := make([]string, 0, len(s))
	for _, p := range s {
		ret = append(ret, p.String())
	}

	return strings.Join(ret, ", ")
}
