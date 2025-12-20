package jsonschema

import (
	"encoding/json/jsontext"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

type Schema struct {
	data
	bytes []byte
}

func NewRefSchema(s string) *Schema {
	return &Schema{data: data{dataCore: dataCore{Ref: s}}}
}

func (s *Schema) Clone() *Schema {
	if s == nil {
		return nil
	}
	return &Schema{data: s.data.Clone(), bytes: s.bytes}
}

func (s *Schema) IsEmpty() bool {
	if s == nil {
		return false
	}
	return s.data.IsEmpty()
}

func (s *Schema) IsEmptyExceptNot() bool {
	if s == nil {
		return false
	}
	return s.data.IsEmptyExceptNot()
}

func (s *Schema) ID() string {
	if s.MetaID != "" {
		return s.MetaID
	}
	if s.ExplicitID != "" {
		return s.ExplicitID
	}
	if s.Ref != "" {
		return s.Ref
	}
	return s.Hash().String()
}

func (s *Schema) Hash() uuid.UUID {
	return util.HashFNV128UUID(util.ToJSONCompact(s))
}

func (s *Schema) String() string {
	ret := fmt.Sprint(s.Type)
	if s.Format != "" {
		ret += "; format=" + s.Format
	}
	return ret
}

func (s *Schema) Summary() string {
	st := s.SchemaType()

	ret := util.NewStringSlice("[" + st.String() + "]")
	if l := s.Properties.Length(); l > 0 {
		ret.Push(util.StringPlural(l, "property"))
	}
	if l := len(s.OneOf); l > 0 {
		ret.Push(util.StringPlural(l, "one-of item"))
	}
	if l := len(s.AnyOf); l > 0 {
		ret.Push(util.StringPlural(l, "any-of item"))
	}
	return ret.Join(", ")
}

func (s *Schema) AddMetadata(k string, v any) {
	if s.Unknown == nil {
		s.Unknown = map[string]jsontext.Value{}
	}
	s.Unknown["x-"+k] = util.ToJSONBytes(v, true)
}

func (s *Schema) GetMetadata() util.ValueMap {
	ret := make(util.ValueMap, len(s.Unknown))
	for k, v := range s.Unknown {
		key := util.Choose(strings.HasPrefix(k, "x-"), k[2:], k)
		ret[key] = util.FromJSONAnyOK(v)
	}
	return ret
}

func (s *Schema) OriginalBytes() []byte {
	if len(s.bytes) == 0 {
		return util.ToJSONBytes(s, true)
	}
	return s.bytes
}

func (s *Schema) IsDeprecated() (bool, string) {
	if s.Deprecated == nil {
		return false, ""
	}
	switch t := s.Deprecated.(type) {
	case bool:
		return t, ""
	case string:
		return true, t
	default:
		return true, fmt.Sprintf("unknown type [%T] for deprecated", s.Deprecated)
	}
}

func (s *Schema) ToFieldDescs() (util.FieldDescs, error) {
	if s.Properties == nil || len(s.Properties.Order) == 0 {
		return nil, errors.New("schema must contain properties")
	}
	ret := make(util.FieldDescs, 0, s.Properties.Length())
	for _, propKey := range s.Properties.Keys() {
		prop := s.Properties.GetSimple(propKey)
		md := s.GetMetadata()
		title := md.GetStringOpt("title")
		typ := fmt.Sprint(prop.Type)
		if prop.Format != "" {
			typ = prop.Format
		}
		d := fmt.Sprint(prop.Default)
		ret = append(ret, &util.FieldDesc{Key: propKey, Title: title, Description: prop.Description, Type: typ, Default: d})
	}
	return ret, nil
}

func (s *Schema) Definitions() *util.OrderedMap[*Schema] {
	if s.Defs.Empty() {
		return s.ExplicitDefs
	}
	if s.ExplicitDefs.Empty() {
		return s.Defs
	}
	return s.ExplicitDefs.Merge(s.Defs)
}
