// Content managed by Project Forge, see [projectforge.md] for details.
package search

import (
	"projectforge.dev/app/lib/filter"
)

type Params struct {
	Q  string          `json:"q"`
	PS filter.ParamSet `json:"ps,omitempty"`
}

func (r *Params) String() string {
	return r.Q
}
