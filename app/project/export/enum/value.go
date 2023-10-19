package enum

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type valueMarshal struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Value struct {
	Key         string
	Title       string
	Description string
	Simple      bool
}

func (x *Value) MarshalJSON() ([]byte, error) {
	if x.Simple {
		return util.ToJSONBytes(x.Key, false), nil
	}
	return util.ToJSONBytes(&valueMarshal{Key: x.Key, Title: x.Title, Description: x.Description}, false), nil
}

func (x *Value) UnmarshalJSON(data []byte) error {
	if strings.Contains(string(data), "{") {
		var v valueMarshal
		if err := util.FromJSON(data, &v); err != nil {
			return err
		}
		x.Key = v.Key
		x.Title = v.Title
		x.Description = v.Description
		x.Simple = false
		return nil
	}
	str := ""
	if err := util.FromJSON(data, &str); err != nil {
		return err
	}
	x.Key = str
	x.Title = ""
	x.Description = ""
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
		return (!x.Simple) || x.Title != "" || x.Description != ""
	})
}

func (v Values) Titles() []string {
	return lo.Map(v, func(x *Value, _ int) string {
		if x.Title != "" {
			return x.Title
		}
		return x.Key
	})
}
