package model

import (
	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app/util"
)

type Models []*Model

func (m Models) Get(pkg util.Pkg, key string) *Model {
	for _, x := range m {
		if x.Pkg.Equals(pkg) && x.Key == key {
			return x
		}
	}
	return nil
}

func (m Models) Sort() {
	slices.SortFunc(m, func(l *Model, r *Model) bool {
		return l.Key < r.Key
	})
}

func (m Models) Names() []string {
	ret := make([]string, 0, len(m))
	for _, md := range m {
		ret = append(ret, md.Key)
	}
	return ret
}
