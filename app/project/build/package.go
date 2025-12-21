package build

import (
	"cmp"
	"runtime"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var ScriptExtension = func() string {
	if runtime.GOOS == OSWindows {
		return "bat"
	}
	return "sh"
}()

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

func (p Pkgs) Sort() Pkgs {
	slices.SortFunc(p, func(l *Pkg, r *Pkg) int {
		return cmp.Compare(strings.ToLower(l.Path), strings.ToLower(r.Path))
	})
	return p
}

func (p Pkgs) ToGraph(prefix string) string {
	if prefix != "" {
		prefix += "/"
	}
	ret := &util.StringSlice{}
	add := func(s string, args ...any) {
		ret.Pushf(s, args...)
	}
	add("graph LR;")
	lo.ForEach(p, func(pkg *Pkg, _ int) {
		lo.ForEach(pkg.Deps, func(d string, _ int) {
			add("\t%s --> %s", strings.TrimPrefix(pkg.Path, prefix), strings.TrimPrefix(d, prefix))
		})
	})
	return ret.String()
}
