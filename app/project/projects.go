package project

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
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
	ret := &util.StringSlice{}
	lo.ForEach(p, func(prj *Project, _ int) {
		lo.ForEach(prj.Modules, func(mod string, _ int) {
			if !lo.Contains(ret.Slice, mod) {
				ret.Push(mod)
			}
		})
	})
	return util.ArraySorted(ret.Slice)
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

func (p Projects) WithModules(modules ...string) Projects {
	return lo.Filter(p, func(x *Project, _ int) bool {
		return len(lo.Intersect(x.Modules, modules)) > 0
	})
}

func (p Projects) WithoutModules(modules ...string) Projects {
	return lo.Filter(p, func(x *Project, _ int) bool {
		return len(lo.Intersect(x.Modules, modules)) == 0
	})
}

func (p Projects) WithTags(tags ...string) Projects {
	return lo.Filter(p, func(x *Project, _ int) bool {
		return len(lo.Intersect(x.Tags, tags)) > 0
	})
}

func (p Projects) WithoutTags(tags ...string) Projects {
	return lo.Filter(p, func(x *Project, _ int) bool {
		return len(lo.Intersect(x.Tags, tags)) == 0
	})
}

func (p Projects) Tags() []string {
	ret := &util.StringSlice{}
	lo.ForEach(p, func(prj *Project, _ int) {
		for _, t := range prj.Tags {
			if !lo.Contains(ret.Slice, t) {
				ret.Push(t)
				break
			}
		}
	})
	return util.ArraySorted(ret.Slice)
}

func (p Projects) ToCSV() ([]string, [][]string) {
	return Fields, lo.Map(p, func(prj *Project, _ int) []string {
		return prj.Strings()
	})
}
