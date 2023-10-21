package enum

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type valueMarshal struct {
	Key         string `json:"key"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

type Value struct {
	Key         string
	Name        string
	Description string
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
	return util.ToJSONBytes(&valueMarshal{Key: x.Key, Name: x.Name, Description: x.Description, Default: x.Default}, false), nil
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
		return (!x.Simple) || x.Name != "" || x.Description != ""
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
