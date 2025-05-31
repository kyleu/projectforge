package jsonschema

import (
	"fmt"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

const KeyExtension = ".schema.json"

func (s *Schema) String() string {
	ret := fmt.Sprint(s.Type)
	if s.Format != "" {
		ret += "; format=" + s.Format
	}
	return ret
}

func (s *Schema) AddMetadata(k string, v any) {
	if s.Metadata == nil {
		s.Metadata = util.ValueMap{}
	}
	s.Metadata[k] = v
}

func (s *Schema) ToFieldDescs() (util.FieldDescs, error) {
	if s.Properties == nil || len(s.Properties.Order) == 0 {
		return nil, errors.New("schema must contain properties")
	}
	ret := make(util.FieldDescs, 0, len(s.Properties.Order))
	for _, propKey := range s.Properties.Order {
		prop := s.Properties.GetSimple(propKey)
		md := prop.Metadata
		if md == nil {
			md = util.ValueMap{}
		}
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
