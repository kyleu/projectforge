package model

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type ImportType string

const (
	ImportTypeInternal ImportType = "internal"
	ImportTypeExternal ImportType = "external"
	ImportTypeApp      ImportType = "app"
)

type Import struct {
	Type      ImportType `json:"type"`
	Value     string     `json:"value"`
	Alias     string     `json:"alias,omitempty"`
	Supported []string   `json:"supported,omitempty"`
}

func NewImport(t ImportType, v string) *Import {
	return &Import{Type: t, Value: v}
}

func (i *Import) WithAlias(a string) *Import {
	i.Alias = a
	return i
}

func (i *Import) Render() string {
	var alias string
	if i.Alias != "" {
		alias = i.Alias + " "
	}
	return fmt.Sprintf("%s%q", alias, i.Value)
}

func (i *Import) Equals(x *Import) bool {
	return x.Type == i.Type && x.Value == i.Value && x.Alias == i.Alias
}

func (i *Import) Supports(key string) bool {
	return len(i.Supported) == 0 || slices.Contains(i.Supported, key)
}

type Imports []*Import

func (i Imports) Render(linebreak string) string {
	if len(i) == 1 {
		return fmt.Sprintf("import %s", i[0].Render())
	}
	ret := util.NewStringSlice([]string{"import ("})
	ret.Push(i.toStrings("\t")...)
	ret.Push(")")
	return ret.Join(linebreak)
}

func (i Imports) RenderHTML(linebreak string) string {
	if len(i) == 1 {
		return fmt.Sprintf("{%% import %s %%}", i[0].Render())
	}
	ret := util.NewStringSlice([]string{"{%% import ("})
	ret.Push(i.toStrings("  ")...)
	ret.Push(") %%}")
	return ret.Join(linebreak)
}

func (i Imports) toStrings(whitespace string) []string {
	ret := &util.StringSlice{}
	add := func(x []string, lf bool) {
		if len(x) > 0 {
			if lf {
				ret.Push("")
			}
			lo.ForEach(x, func(item string, _ int) {
				ret.Push(whitespace + item)
			})
		}
	}

	internal := i.renderByType(ImportTypeInternal)
	external := i.renderByType(ImportTypeExternal)
	app := i.renderByType(ImportTypeApp)
	add(internal, false)
	add(external, len(internal) > 0)
	add(app, len(external) > 0 || len(internal) > 0)

	return ret.Slice
}

func (i Imports) renderByType(t ImportType) []string {
	ret := lo.FilterMap(i, func(x *Import, _ int) (string, bool) {
		return x.Render(), x.Type == t
	})
	slices.SortFunc(ret, func(l string, r string) int {
		if lIdx := strings.Index(l, " "); lIdx > -1 {
			l = l[lIdx+1:]
		}
		if rIdx := strings.Index(r, " "); rIdx > -1 {
			r = r[rIdx+1:]
		}
		return cmp.Compare(l, r)
	})
	return ret
}

func (i Imports) Add(imports ...*Import) Imports {
	if i == nil {
		return slices.Clone(imports)
	}
	ret := slices.Clone(i)
	lo.ForEach(imports, func(imp *Import, _ int) {
		hit := lo.ContainsBy(ret, func(x *Import) bool {
			return x.Value == imp.Value
		})
		if !hit {
			ret = append(ret, imp)
		}
	})
	return ret
}

func (i Imports) Supporting(key string) Imports {
	return lo.Filter(i, func(x *Import, _ int) bool {
		return x.Supports(key)
	})
}
