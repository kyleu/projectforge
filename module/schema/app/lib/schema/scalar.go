package schema

import (
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Scalar struct {
	Pkg         util.Pkg  `json:"pkg,omitempty"`
	Key         string    `json:"key"`
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}

type Scalars []*Scalar

func (s Scalars) Get(pkg util.Pkg, key string) *Scalar {
	return lo.FindOrElse(s, nil, func(x *Scalar) bool {
		return x.Pkg.Equals(pkg) && x.Key == key
	})
}
