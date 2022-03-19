package model

import (
	"{{{ .Package }}}/app/util"
)

type Package struct {
	Key           string   `json:"key"`
	Title         string   `json:"title,omitempty"`
	Pkg           util.Pkg `json:"pkg"`
	ChildModels   Models   `json:"childModels,omitempty"`
	ChildPackages Packages `json:"childPackages,omitempty"`
}

func (p *Package) Name() string {
	if p.Title == "" {
		return p.Key
	}
	return p.Title
}

func (p *Package) Path() util.Pkg {
	return p.Pkg.With(p.Key)
}

func (p *Package) PathString() string {
	return p.Pkg.ToPath(p.Key)
}

func (p *Package) GetPkg(key string, createIfMissing bool) *Package {
	for _, x := range p.ChildPackages {
		if x.Key == key {
			return x
		}
	}
	if createIfMissing {
		pkgs := p.Path()
		if pkgs[0] == "_root" {
			pkgs = pkgs[1:]
		}
		x := &Package{Key: key, Pkg: pkgs}
		p.ChildPackages = append(p.ChildPackages, x)
		return x
	}
	return nil
}

func (p *Package) GetModel(key string) *Model {
	for _, x := range p.ChildModels {
		if x.Key == key {
			return x
		}
	}
	return nil
}

func (p *Package) Add(pkg util.Pkg, m *Model) {
	if len(pkg) == 0 {
		p.ChildModels = append(p.ChildModels, m)
	} else {
		x := p.GetPkg(pkg[0], true)
		x.Add(pkg[1:], m)
	}
}

func (p *Package) Get(paths []string) (any, []string) {
	if len(paths) == 0 {
		return p, nil
	}
	if x := p.GetPkg(paths[0], false); x != nil {
		return x.Get(paths[1:])
	}
	m := p.GetModel(paths[0])
	if m == nil {
		return p, paths
	}
	return m, paths[1:]
}

type Packages []*Package

func ToModelPackage(models Models) *Package {
	ret := &Package{Key: "_root", Title: "Project Root"}
	for _, m := range models {
		ret.Add(m.Pkg, m)
	}
	return ret
}
