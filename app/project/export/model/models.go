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

func (m Models) Sorted() Models {
	ret := make(Models, 0, len(m))
	for _, mdl := range m {
		if curr := ret.Get(mdl.Name); curr == nil {
			for _, n := range m.withDeps(mdl) {
				if x := ret.Get(n.Name); x == nil {
					ret = append(ret, n)
				}
			}
		}
	}
	return ret
}

func (m Models) withDeps(mdl *Model) Models {
	var deps Models
	for _, rel := range mdl.Relations {
		if deps.Get(rel.Table) == nil {
			deps = append(deps, m.withDeps(m.Get(rel.Table))...)
		}
	}
	if deps.Get(mdl.Name) == nil {
		deps = append(deps, mdl)
	}
	return deps
}

func (m Models) HasSearch() bool {
	for _, x := range m {
		if len(x.AllSearches()) > 0 {
			return true
		}
	}
	return false
}

func (m Models) HasSeedData() bool {
	for _, x := range m {
		if len(x.SeedData) > 0 {
			return true
		}
	}
	return false
}

func (m Models) ForGroup(pth ...string) Models {
	var ret Models
	for _, x := range m {
		if slices.Equal(x.Group, pth) {
			ret = append(ret, x)
		}
	}
	return ret
}
