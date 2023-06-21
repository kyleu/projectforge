package golang

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type ImportType string

const (
	ImportTypeInternal ImportType = "internal"
	ImportTypeExternal ImportType = "external"
	ImportTypeApp      ImportType = "app"
)

type Import struct {
	Type  ImportType
	Value string
	Alias string
}

func (i *Import) WithAlias(a string) *Import {
	i.Alias = a
	return i
}

func (i *Import) Render() string {
	alias := ""
	if i.Alias != "" {
		alias = i.Alias + " "
	}
	return fmt.Sprintf("%s%q", alias, i.Value)
}

func (i *Import) Equals(x *Import) bool {
	return x.Type == i.Type && x.Value == i.Value && x.Alias == i.Alias
}

func NewImport(t ImportType, v string) *Import {
	return &Import{Type: t, Value: v}
}

type Imports []*Import

func (i Imports) Render() string {
	if len(i) == 1 {
		return fmt.Sprintf("import %s", i[0].Render())
	}
	ret := []string{"import ("}

	add := func(x []string, lf bool) {
		if len(x) > 0 {
			if lf {
				ret = append(ret, "")
			}
			lo.ForEach(x, func(item string, _ int) {
				ret = append(ret, fmt.Sprintf("\t%s", item))
			})
		}
	}

	internal := i.ByType(ImportTypeInternal)
	external := i.ByType(ImportTypeExternal)
	app := i.ByType(ImportTypeApp)
	add(internal, false)
	add(external, len(internal) > 0)
	add(app, len(external) > 0 || len(internal) > 0)

	ret = append(ret, ")")
	return strings.Join(ret, "\n")
}

func (i Imports) RenderHTML() string {
	if len(i) == 1 {
		return fmt.Sprintf("{%% import %s %%}", i[0].Render())
	}
	ret := []string{"{%% import ("}

	add := func(x []string, lf bool) {
		if len(x) > 0 {
			if lf {
				ret = append(ret, "")
			}
			lo.ForEach(x, func(item string, _ int) {
				ret = append(ret, fmt.Sprintf("  %s", item))
			})
		}
	}

	internal := i.ByType(ImportTypeInternal)
	external := i.ByType(ImportTypeExternal)
	app := i.ByType(ImportTypeApp)
	add(internal, false)
	add(external, len(internal) > 0)
	add(app, len(external) > 0 || len(internal) > 0)

	ret = append(ret, ") %%}")
	return strings.Join(ret, "\n")
}

func (i Imports) ByType(t ImportType) []string {
	ret := lo.FilterMap(i, func(x *Import, _ int) (string, bool) {
		return x.Render(), x.Type == t
	})
	slices.Sort(ret)
	return ret
}

func (i Imports) Add(imports ...*Import) Imports {
	if i == nil {
		return append(Imports{}, imports...)
	}
	lo.ForEach(imports, func(imp *Import, _ int) {
		hit := lo.ContainsBy(i, func(x *Import) bool {
			return x.Value == imp.Value
		})
		if !hit {
			i = append(i, imp)
		}
	})
	return i
}
