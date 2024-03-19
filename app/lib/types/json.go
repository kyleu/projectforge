// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import "projectforge.dev/projectforge/app/util"

const KeyJSON = "json"

type JSON struct {
	IsObject bool `json:"obj,omitempty"`
	IsArray  bool `json:"arr,omitempty"`
}

var _ Type = (*JSON)(nil)

func (x *JSON) Key() string {
	return KeyJSON
}

func (x *JSON) Sortable() bool {
	return false
}

func (x *JSON) Scalar() bool {
	return false
}

func (x *JSON) String() string {
	return x.Key()
}

func (x *JSON) From(v any) any {
	if s, ok := v.(string); ok {
		ret, _ := util.FromJSONAny([]byte(s))
		return ret
	}
	return invalidInput(x.Key(), v)
}

func (x *JSON) Default(string) any {
	return "{}"
}

func NewJSON() *Wrapped {
	return Wrap(&JSON{})
}

func NewJSONArgs(obj bool, arr bool) *Wrapped {
	return Wrap(&JSON{IsObject: obj, IsArray: arr})
}
