package build

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Pkg struct {
	Path  string   `json:"path"`
	Files []string `json:"files,omitempty"`
	Deps  []string `json:"deps,omitempty"`
}

func (p *Pkg) AddDep(s string) {
	for _, x := range p.Deps {
		if x == s {
			return
		}
	}
	p.Deps = append(p.Deps, s)
}

type Pkgs []*Pkg

func (p Pkgs) Get(s string) *Pkg {
	for _, x := range p {
		if x.Path == s {
			return x
		}
	}
	return nil
}

func (p Pkgs) Sort() {
	slices.SortFunc(p, func(l *Pkg, r *Pkg) bool {
		return l.Path < r.Path
	})
}

func (p Pkgs) ToGraph(prefix string) string {
	if prefix != "" {
		prefix += "/"
	}
	var ret []string
	add := func(s string, args ...any) {
		ret = append(ret, fmt.Sprintf(s, args...))
	}
	add("graph LR;")
	for _, pkg := range p {
		for _, d := range pkg.Deps {
			add("\t%s --> %s", strings.TrimPrefix(pkg.Path, prefix), strings.TrimPrefix(d, prefix))
		}
	}
	return strings.Join(ret, "\n")
}
