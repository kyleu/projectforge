package project

import (
	"golang.org/x/exp/slices"
)

type Projects []*Project

func (p Projects) Get(key string) *Project {
	for _, x := range p {
		if x.Key == key {
			return x
		}
	}
	return nil
}

func (p Projects) Root() *Project {
	for _, x := range p {
		if x.Path == "." {
			return x
		}
	}
	return nil
}

func (p Projects) AllModules() []string {
	var ret []string
	for _, prj := range p {
		for _, mod := range prj.Modules {
			if !slices.Contains(ret, mod) {
				ret = append(ret, mod)
			}
		}
	}
	slices.Sort(ret)
	return ret
}

func (p Projects) Keys() []string {
	ret := make([]string, 0, len(p))
	for _, prj := range p {
		ret = append(ret, prj.Key)
	}
	return ret
}

func (p Projects) Titles() []string {
	ret := make([]string, 0, len(p))
	for _, prj := range p {
		ret = append(ret, prj.Title())
	}
	return ret
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
	for _, prj := range p {
		var hit bool
		for _, t := range tags {
			if slices.Contains(prj.Tags, t) {
				hit = true
				break
			}
		}
		if !hit {
			ret = append(ret, prj)
		}
	}
	return ret
}

func (p Projects) Tags() []string {
	var ret []string
	for _, prj := range p {
		for _, t := range prj.Tags {
			if !slices.Contains(ret, t) {
				ret = append(ret, t)
				break
			}
		}
	}
	slices.Sort(ret)
	return ret
}
