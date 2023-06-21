package project

import (
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type Projects []*Project

func (p Projects) Get(key string) *Project {
	return lo.FindOrElse(p, nil, func(x *Project) bool {
		return x.Key == key
	})
}

func (p Projects) Root() *Project {
	return lo.FindOrElse(p, nil, func(x *Project) bool {
		return x.Path == "."
	})
}

func (p Projects) AllModules() []string {
	var ret []string
	lo.ForEach(p, func(prj *Project, _ int) {
		lo.ForEach(prj.Modules, func(mod string, _ int) {
			if !slices.Contains(ret, mod) {
				ret = append(ret, mod)
			}
		})
	})
	slices.Sort(ret)
	return ret
}

func (p Projects) Keys() []string {
	return lo.Map(p, func(prj *Project, _ int) string {
		return prj.Key
	})
}

func (p Projects) Titles() []string {
	return lo.Map(p, func(prj *Project, _ int) string {
		return prj.Title()
	})
}

func (p Projects) WithTags(tags ...string) Projects {
	ret := make(Projects, 0, len(p))
	for _, prj := range p {
		for _, t := range tags {
			if slices.Contains(prj.Tags, t) {
				ret = append(ret, prj)
				break
			}
		}
	}
	return ret
}

func (p Projects) WithoutTags(tags ...string) Projects {
	ret := make(Projects, 0, len(p))
	lo.ForEach(p, func(prj *Project, _ int) {
		hit := lo.ContainsBy(tags, func(t string) bool {
			return slices.Contains(prj.Tags, t)
		})
		if !hit {
			ret = append(ret, prj)
		}
	})
	return ret
}

func (p Projects) Tags() []string {
	var ret []string
	lo.ForEach(p, func(prj *Project, _ int) {
		for _, t := range prj.Tags {
			if !slices.Contains(ret, t) {
				ret = append(ret, t)
				break
			}
		}
	})
	slices.Sort(ret)
	return ret
}
