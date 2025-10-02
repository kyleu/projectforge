package enum

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/types"
	"{{{ .Package }}}/app/util"
)

const defaultIcon = "hammer"

type Enum struct {
	Name           string        `json:"name"`
	Package        string        `json:"package"`
	Group          []string      `json:"group,omitempty"`
	Schema         string        `json:"schema,omitzero"`
	Description    string        `json:"description,omitzero"`
	Icon           string        `json:"icon,omitzero"`
	Values         Values        `json:"values,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
	TitleOverride  string        `json:"title,omitzero"`
	ProperOverride string        `json:"proper,omitzero"`
	RouteOverride  string        `json:"route,omitzero"`
	Config         util.ValueMap `json:"config,omitzero"`
	Acronyms       []string      `json:"-"`
}

func (e *Enum) Title() string {
	if e.TitleOverride == "" {
		return util.StringToTitle(e.Name, e.Acronyms...)
	}
	return e.TitleOverride
}

func (e *Enum) TitleLower() string {
	return strings.ToLower(e.Title())
}

func (e *Enum) Proper() string {
	if e.ProperOverride == "" {
		return util.StringToProper(e.Name, e.Acronyms...)
	}
	return util.StringToProper(e.ProperOverride, e.Acronyms...)
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
	return util.StringToCamel(e.Name, e.Acronyms...)
}

func (e *Enum) CamelLower() string {
	return strings.ToLower(e.Camel())
}

func (e *Enum) Kebab() string {
	return util.StringToKebab(e.Name, e.Acronyms...)
}

func (e *Enum) ExtraFields() *util.OrderedMap[string] {
	ret := util.NewOrderedMap[string](false, 0)
	for _, v := range e.Values {
		if v.Extra == nil {
			continue
		}
		for _, k := range v.Extra.Order {
			x := v.Extra.GetSimple(k)
			if _, exists := ret.Get(k); exists {
				continue
			}
			var typ string
			switch x.(type) {
			case string:
				typ = types.KeyString
			case float64:
				typ = types.KeyFloat
			case int, int32, int64:
				typ = types.KeyInt
			case bool:
				typ = types.KeyBool
			}
			if x := e.Config.GetStringOpt("type:" + k); x != "" {
				switch x {
				case types.KeyBool, types.KeyFloat, types.KeyInt, types.KeyString, types.KeyTimestamp:
					typ = x
				default:
					typ = "unknown config type [" + x + "]"
				}
			}
			ret.Set(k, typ)
		}
	}
	return ret
}

func (e *Enum) ExtraFieldValues(k string) ([]any, bool) {
	ret := make([]any, 0, len(e.Values))
	for _, v := range e.Values {
		if v.Extra == nil {
			continue
		}
		if x, ok := v.Extra.Get(k); ok && x != nil {
			ret = append(ret, x)
		}
	}
	return ret, len(lo.Uniq(ret)) == len(ret)
}

func (e *Enum) ID() string {
	return util.StringPath(e.PackageWithGroup(""), e.Name)
}

func (e *Enum) PackageWithGroup(prefix string) string {
	if x := e.Config.GetStringOpt("pkg-" + prefix); x != "" {
		return x
	}
	if len(e.Group) == 0 {
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
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
	return util.StringJoin(x, "/")
}

func (e *Enum) HasTag(t string) bool {
	return lo.Contains(e.Tags, t)
}

func (e *Enum) HasValueIcons() bool {
	return lo.ContainsBy(e.Values, func(x *Value) bool {
		return x.Icon != ""
	})
}

func (e *Enum) Breadcrumbs() string {
	ret := lo.Map(e.Group, func(g string, _ int) string {
		return fmt.Sprintf("%q", g)
	})
	ret = append(ret, fmt.Sprintf("%q", e.Package))
	return util.StringJoin(ret, ", ")
}

func (e *Enum) ValuesCamel() []string {
	return lo.Map(e.Values, func(x *Value, _ int) string {
		return util.StringToProper(x.Key, e.Acronyms...)
	})
}

func (e *Enum) Simple() bool {
	return e.Values.AllSimple()
}

func (e *Enum) Clone() *Enum {
	return &Enum{
		Name:           e.Name,
		Package:        e.Package,
		Group:          e.Group,
		Schema:         e.Schema,
		Description:    e.Description,
		Icon:           e.Icon,
		Values:         e.Values.Clone(),
		Tags:           util.ArrayCopy(e.Tags),
		TitleOverride:  e.TitleOverride,
		ProperOverride: e.ProperOverride,
		RouteOverride:  e.RouteOverride,
		Config:         e.Config.Clone(),
		Acronyms:       util.ArrayCopy(e.Acronyms),
	}
}

func (e *Enum) AddTag(t string) {
	if !lo.Contains(e.Tags, t) {
		e.Tags = util.ArraySorted(append(e.Tags, t))
	}
}

type Enums []*Enum

func (e Enums) Get(key string) *Enum {
	return lo.FindOrElse(e, nil, func(x *Enum) bool {
		return x.Name == key
	})
}

func (e Enums) Clone() Enums {
	return lo.Map(e, func(x *Enum, index int) *Enum {
		return x.Clone()
	})
}

func (e Enums) WithTag(tag string) Enums {
	return lo.Filter(e, func(x *Enum, _ int) bool {
		return x.HasTag(tag)
	})
}

func (e Enums) WithoutTag(tag string) Enums {
	return lo.Reject(e, func(x *Enum, _ int) bool {
		return x.HasTag(tag)
	})
}
