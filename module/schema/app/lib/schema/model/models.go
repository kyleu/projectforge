package model

import (
	"cmp"
	"slices"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Models []*Model

func (m Models) Get(pkg util.Pkg, key string) *Model {
	return lo.FindOrElse(m, nil, func(x *Model) bool {
		return x.Pkg.Equals(pkg) && x.Key == key
	})
}

func (m Models) Sort() {
	slices.SortFunc(m, func(l *Model, r *Model) int {
		return cmp.Compare(l.Key, r.Key)
	})
}

func (m Models) Names() []string {
	return lo.Map(m, func(x *Model, _ int) string {
		return x.Key
	})
}
