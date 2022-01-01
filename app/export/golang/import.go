package golang

import (
	"fmt"
	"sort"
	"strings"
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
}

func NewImport(t ImportType, v string) *Import {
	return &Import{Type: t, Value: v}
}

type Imports []*Import

func (i Imports) Render() string {
	if len(i) == 1 {
		return fmt.Sprintf("import %q", i[0].Value)
	}
	ret := []string{"import ("}

	add := func(x []string, lf bool) {
		if len(x) > 0 {
			if lf {
				ret = append(ret, "")
			}
			for _, x := range x {
				ret = append(ret, fmt.Sprintf("\t%q", x))
			}
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
		return fmt.Sprintf("{%% import %q %%}", i[0].Value)
	}
	ret := []string{"{%% import ("}

	add := func(x []string, lf bool) {
		if len(x) > 0 {
			if lf {
				ret = append(ret, "")
			}
			for _, x := range x {
				ret = append(ret, fmt.Sprintf("  %q", x))
			}
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
	var ret []string
	for _, x := range i {
		if x.Type == t {
			ret = append(ret, x.Value)
		}
	}
	sort.Strings(ret)
	return ret
}

func (i Imports) Add(imports ...*Import) Imports {
	for _, imp := range imports {
		hit := false
		for _, x := range i {
			if x.Value == imp.Value {
				hit = true
				break
			}
		}
		if !hit {
			i = append(i, imp)
		}
	}
	return i
}
