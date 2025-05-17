package enum

import (
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

var ValueFieldDescs = util.FieldDescs{
	{Key: "key", Title: "Key", Description: "The key of the enum"},
	{Key: "name", Title: "Name", Description: "The name of the enum"},
	{Key: "description", Title: "Description", Description: "The description of the enum"},
	{Key: "icon", Title: "Icon", Description: "The icon of the enum", Type: "icon"},
	{Key: "default", Title: "Default", Description: "Indicates if this is the default value", Type: "bool"},
	// {Key: "extra", Title: "Extra", Description: "The X of the column"},
}

type valueMarshal struct {
	Key         string                `json:"key"`
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	Icon        string                `json:"icon,omitempty"`
	Default     bool                  `json:"default,omitempty"`
	Extra       *util.OrderedMap[any] `json:"extra,omitempty"`
}

type Value struct {
	Key         string
	Name        string
	Description string
	Icon        string
	Extra       *util.OrderedMap[any]
	Default     bool
	Simple      bool
}

func (x *Value) Proper(acronyms ...string) string {
	return util.StringToProper(x.Key, acronyms...)
}

func (x *Value) Title() string {
	return util.StringToTitle(x.Key)
}

func (x *Value) ToOrderedMap(includeKey bool) *util.OrderedMap[any] {
	ret := util.NewOrderedMap[any](false, 6)
	if includeKey {
		ret.Set("key", x.Key)
	}
	if x.Name != "" {
		ret.Set("name", x.Name)
	}
	if x.Description != "" {
		ret.Set("description", x.Description)
	}
	if x.Icon != "" {
		ret.Set("icon", x.Icon)
	}
	if x.Default {
		ret.Set("default", x.Default)
	}
	if x.Extra != nil && len(x.Extra.Order) > 0 {
		ret.Set("extra", x.Extra)
	}
	return ret
}

func (x *Value) MarshalJSON() ([]byte, error) {
	if x.Simple {
		return util.ToJSONBytes(x.Key, false), nil
	}
	marshaller := &valueMarshal{Key: x.Key, Name: x.Name, Description: x.Description, Icon: x.Icon, Extra: x.Extra, Default: x.Default}
	return util.ToJSONBytes(marshaller, false), nil
}

func (x *Value) UnmarshalJSON(data []byte) error {
	if strings.Contains(string(data), "{") {
		var v valueMarshal
		if err := util.FromJSON(data, &v); err != nil {
			return err
		}
		x.Key = v.Key
		x.Name = v.Name
		x.Description = v.Description
		x.Icon = v.Icon
		x.Extra = v.Extra
		x.Default = v.Default
		x.Simple = false
		return nil
	}
	str, err := util.FromJSONString(data)
	if err != nil {
		return err
	}
	x.Key = str
	x.Simple = true
	return nil
}

func (x *Value) Clone() *Value {
	return &Value{Key: x.Key, Name: x.Name, Description: x.Description, Icon: x.Icon, Extra: x.Extra.Clone(), Default: x.Default, Simple: x.Simple}
}

type Values []*Value

func (v Values) Keys() []string {
	return lo.Map(v, func(x *Value, _ int) string {
		return x.Key
	})
}

func (v Values) Names() []string {
	return lo.Map(v, func(x *Value, _ int) string {
		return x.Name
	})
}

func (v Values) AllSimple() bool {
	return lo.NoneBy(v, func(x *Value) bool {
		return (!x.Simple) || x.Name != "" || x.Description != "" || x.Icon != "" || (x.Extra != nil && len(x.Extra.Map) > 0)
	})
}

func (v Values) Default() *Value {
	return lo.FindOrElse(v, nil, func(x *Value) bool {
		return x.Default
	})
}

func (v Values) Titles() []string {
	return lo.Map(v, func(x *Value, _ int) string {
		if x.Name != "" {
			return x.Title()
		}
		return x.Key
	})
}

func (v Values) Clone() Values {
	return lo.Map(v, func(x *Value, index int) *Value {
		return x.Clone()
	})
}
