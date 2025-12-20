package typescript

import (
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type TSImport struct {
	Ref    string `json:"ref,omitzero"`
	Path   string `json:"path,omitzero"`
	Prefix string `json:"prefix,omitzero"`
}

func newImport(ref string, path string, prefix ...string) *TSImport {
	return &TSImport{Ref: ref, Path: path, Prefix: strings.Join(prefix, "")}
}

func (t *TSImport) RefString() string {
	return util.Choose(t.Prefix == "", "", t.Prefix+" ") + strings.TrimPrefix(t.Ref, "*")
}

type TSImports []*TSImport

func (t TSImports) With(xx ...*TSImport) TSImports {
	ret := t
	for _, x := range xx {
		ret = ret.WithOne(x)
	}
	return ret
}

func (t TSImports) WithOne(x *TSImport) TSImports {
	for _, imp := range t {
		if imp.Ref == x.Ref && imp.Path == x.Path {
			if x.Prefix != "" {
				imp.Prefix = x.Prefix
			}
			return t
		}
	}
	return append(t, x)
}

func (t TSImports) RefStrings() []string {
	return lo.Map(t, func(x *TSImport, _ int) string {
		return x.RefString()
	})
}

func (t TSImports) Paths() []string {
	ret := lo.Uniq(lo.Map(t, func(x *TSImport, _ int) string {
		return x.Path
	}))
	slices.Sort(ret)
	return ret
}

func (t TSImports) Strings() []string {
	ret := util.NewStringSliceWithSize(len(t.Paths()))
	for _, pth := range t.Paths() {
		matches := t.ByPath(pth)
		ret.Pushf("import { %s } from %q;", strings.Join(matches.RefStrings(), ", "), pth)
	}
	ret.Sort()
	return ret.Slice
}

func (t TSImports) ByPath(x string) TSImports {
	return lo.Filter(t, func(i *TSImport, _ int) bool {
		return i.Path == x
	})
}
