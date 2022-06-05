package model

import (
	"golang.org/x/exp/slices"
)

type Models []*Model

func (m Models) Get(n string) *Model {
	for _, x := range m {
		if x.Name == n {
			return x
		}
	}
	return nil
}

func (m Models) ReverseRelations(t string) Relations {
	var rels Relations
	for _, x := range m {
		for _, rel := range x.Relations {
			if rel.Table == t {
				rels = append(rels, rel.Reverse(x.Name))
			}
		}
	}
	return rels
}

func (m Models) Replace(mdl *Model) Models {
	for idx, curr := range m {
		if curr.Name == mdl.Name {
			m[idx] = mdl
			m.Sort()
			return m
		}
	}
	m = append(m, mdl)
	m.Sort()
	return m
}

func (m Models) Sort() {
	slices.SortFunc(m, func(l *Model, r *Model) bool {
		return l.Offset < r.Offset
	})
}
