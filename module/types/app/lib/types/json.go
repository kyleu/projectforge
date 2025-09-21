package types

import "{{{ .Package }}}/app/util"

const KeyJSON, objStr = util.KeyJSON, "{}"

type JSON struct {
	IsObject bool `json:"obj,omitzero"`
	IsArray  bool `json:"arr,omitzero"`
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
	if s, err := util.Cast[string](v); err == nil {
		ret, _ := util.FromJSONAny([]byte(s))
		return ret
	}
	return invalidInput(x.Key(), v)
}

func (x *JSON) Default(string) any {
	return objStr
}

func NewJSON() *Wrapped {
	return Wrap(&JSON{})
}

func NewJSONArgs(obj bool, arr bool) *Wrapped {
	return Wrap(&JSON{IsObject: obj, IsArray: arr})
}
