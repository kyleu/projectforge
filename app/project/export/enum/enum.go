package enum

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const defaultIcon = "hammer"

type Enum struct {
	Name           string        `json:"name"`
	Package        string        `json:"package"`
	Group          []string      `json:"group,omitempty"`
	Description    string        `json:"description,omitempty"`
	Icon           string        `json:"icon,omitempty"`
	Values         Values        `json:"values,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
	TitleOverride  string        `json:"title,omitempty"`
	ProperOverride string        `json:"proper,omitempty"`
	RouteOverride  string        `json:"route,omitempty"`
	Config         util.ValueMap `json:"config,omitempty"`
}

func (e *Enum) Title() string {
	if e.TitleOverride == "" {
		return e.Proper()
	}
	return e.TitleOverride
}

func (e *Enum) Proper() string {
	if e.ProperOverride == "" {
		return util.StringToCamel(e.Name)
	}
	return util.StringToCamel(e.ProperOverride)
}

func (e *Enum) ProperPlural() string {
	ret := util.StringToPlural(e.Proper())
	if ret == e.Proper() {
		return ret + "Set"
	}
	return ret
}

func (e *Enum) FirstLetter() any {
	return strings.ToLower(e.Name[0:1])
}

func (e *Enum) IconSafe() string {
	if _, ok := util.SVGLibrary[e.Icon]; ok {
		return e.Icon
	}
	return defaultIcon
}

func (e *Enum) Camel() string {
	return util.StringToLowerCamel(e.Name)
}

func (e *Enum) ExtraFields() map[string]string {
	ret := map[string]string{}
	for _, v := range e.Values {
		for k, v := range v.Extra {
			typ := ""
			switch v.(type) {
			case string:
				typ = "string"
			case float64:
				typ = "float"
			case int, int32, int64:
				typ = "int"
			}
			ret[k] = typ
		}
	}
	return ret
}

func (e *Enum) PackageWithGroup(prefix string) string {
	if len(e.Group) == 0 {
		if len(prefix) > 0 && !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		return prefix + e.Package
	}
	var x []string
	if prefix != "" {
		x = append(x, prefix)
	}
	x = append(x, e.Group...)
	x = append(x, e.Package)
	return strings.Join(x, "/")
}

func (e *Enum) HasTag(t string) bool {
	return lo.Contains(e.Tags, t)
}

func (e *Enum) Breadcrumbs() string {
	ret := lo.Map(e.Group, func(g string, _ int) string {
		return fmt.Sprintf("%q", g)
	})
	ret = append(ret, fmt.Sprintf("%q", e.Package))
	return strings.Join(ret, ", ")
}

func (e *Enum) ValuesCamel() []string {
	return lo.Map(e.Values, func(x *Value, _ int) string {
		return util.StringToCamel(x.Key)
	})
}

func (e *Enum) Simple() bool {
	return e.Values.AllSimple()
}

type Enums []*Enum

func (e Enums) Get(key string) *Enum {
	return lo.FindOrElse(e, nil, func(x *Enum) bool {
		return x.Name == key
	})
}
