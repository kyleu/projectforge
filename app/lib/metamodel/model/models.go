package model

import (
	"cmp"
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Models []*Model

func (m Models) Get(n string) *Model {
	return lo.FindOrElse(m, nil, func(x *Model) bool {
		return x.Name == n
	})
}

func (m Models) ReverseRelations(t string) Relations {
	return lo.FlatMap(m, func(x *Model, _ int) []*Relation {
		return lo.FilterMap(x.Relations, func(rel *Relation, _ int) (*Relation, bool) {
			if rel.Table == t {
				return rel.Reverse(x.Name), true
			}
			return nil, false
		})
	})
}

func (m Models) Sorted() Models {
	ret := make(Models, 0, len(m))
	lo.ForEach(m, func(mdl *Model, _ int) {
		if curr := ret.Get(mdl.Name); curr == nil {
			lo.ForEach(m.withDeps(mdl), func(n *Model, _ int) {
				if x := ret.Get(n.Name); x == nil {
					ret = append(ret, n)
				}
			})
		}
	})
	return ret
}

func (m Models) SortedDisplay() Models {
	ret := util.ArrayCopy(m)
	slices.SortFunc(ret, func(l *Model, r *Model) int {
		if l.SortIndex == r.SortIndex {
			return cmp.Compare(l.Name, r.Name)
		}
		return cmp.Compare(l.SortIndex, r.SortIndex)
	})
	return ret
}

func (m Models) SortedRoutes() Models {
	ret := util.ArrayCopy(m)
	slices.SortFunc(ret, func(l *Model, r *Model) int {
		return cmp.Compare(r.PackageWithGroup(""), l.PackageWithGroup(""))
	})
	return ret
}

func (m Models) withDeps(mdl *Model) Models {
	var deps Models
	lo.ForEach(mdl.Relations, func(rel *Relation, _ int) {
		if deps.Get(rel.Table) == nil && rel.Table != mdl.Name {
			deps = append(deps, m.withDeps(m.Get(rel.Table))...)
		}
	})
	if deps.Get(mdl.Name) == nil {
		deps = append(deps, mdl)
	}
	return deps
}

func (m Models) HasSearch() bool {
	return lo.ContainsBy(m, func(x *Model) bool {
		return x.HasSearches()
	})
}

func (m Models) HasSeedData() bool {
	return lo.ContainsBy(m, func(x *Model) bool {
		return len(x.SeedData) > 0
	})
}

func (m Models) ForGroup(pth ...string) Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return slices.Equal(x.Group, pth)
	})
}

func (m Models) Validate(mods []string, groups Groups) error {
	names := util.ValueMap{}
	for _, x := range m {
		if _, ok := names[x.Name]; ok {
			return errors.Errorf("multiple models found with name [%s]", x.Name)
		}
		if err := x.Validate(mods, m, groups); err != nil {
			return err
		}
	}
	return nil
}

func (m Models) WithTag(tag string) Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return x.HasTag(tag)
	})
}

func (m Models) WithoutTag(tag string) Models {
	return lo.Reject(m, func(x *Model, _ int) bool {
		return x.HasTag(tag)
	})
}
