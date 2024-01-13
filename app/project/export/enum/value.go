package enum

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

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

func (x *Value) Proper() string {
	return util.StringToCamel(x.Key)
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
		if x.Extra == nil {
			x.Extra = util.NewOrderedMap[any](false, 0)
		}
		x.Default = v.Default
		x.Simple = false
		return nil
	}
	str := ""
	if err := util.FromJSON(data, &str); err != nil {
		return err
	}
	x.Key = str
	x.Name = ""
	x.Description = ""
	x.Icon = ""
	x.Extra = util.NewOrderedMap[any](false, 0)
	x.Default = false
	x.Simple = true
	return nil
}

type Values []*Value

func (v Values) Keys() []string {
	return lo.Map(v, func(x *Value, _ int) string {
		return x.Key
	})
}

func (v Values) AllSimple() bool {
	return !lo.ContainsBy(v, func(x *Value) bool {
		return (!x.Simple) || x.Name != "" || x.Description != "" || x.Icon != "" || len(x.Extra.Map) > 0
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
			return x.Name
		}
		return x.Key
	})
}
