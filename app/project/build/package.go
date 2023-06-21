package build

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type Pkg struct {
	Path  string   `json:"path"`
	Files []string `json:"files,omitempty"`
	Deps  []string `json:"deps,omitempty"`
}

func (p *Pkg) AddDep(s string) {
	if lo.Contains(p.Deps, s) {
		return
	}
	p.Deps = append(p.Deps, s)
}

type Pkgs []*Pkg

func (p Pkgs) Get(s string) *Pkg {
	return lo.FindOrElse(p, nil, func(x *Pkg) bool {
		return x.Path == s
	})
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
	lo.ForEach(p, func(pkg *Pkg, _ int) {
		lo.ForEach(pkg.Deps, func(d string, index int) {
			add("\t%s --> %s", strings.TrimPrefix(pkg.Path, prefix), strings.TrimPrefix(d, prefix))
		})
	})
	return strings.Join(ret, "\n")
}
