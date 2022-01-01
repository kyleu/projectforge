package search

import (
	"{{{ .Package }}}/app/lib/filter"
)

type Params struct {
	Q  string          `json:"q"`
	PS filter.ParamSet `json:"ps,omitempty"`
}

func (r *Params) String() string {
	return r.Q
}
