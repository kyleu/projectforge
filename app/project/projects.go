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

func (p Projects) WithTags(tags ...string) Projects {
	ret := make(Projects, 0, len(p))
	for _, prj := range p {
		for _, t := range tags {
			if lo.Contains(prj.Tags, t) {
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
			return lo.Contains(prj.Tags, t)
		})
		if !hit {
			ret = append(ret, prj)
		}
	})
	return ret
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
