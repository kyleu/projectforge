// $PF_IGNORE$
package graph

import (
	"{{{ .Package }}}/app/util"
)

type Resolver struct {
	logger util.Logger
}

func NewResolver(logger util.Logger) *Resolver {
	return &Resolver{logger}
}
