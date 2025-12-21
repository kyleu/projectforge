package jsonschema

import (
	"bytes"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func (s *Schema) UnmarshalJSON(data []byte) error {
	str, err := schemaFromJSON(data)
	if err != nil {
		return err
	}
	if str != nil {
		*s = *str
	}
	return nil
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	if s.IsEmpty() {
		return []byte("true"), nil
	}
	if s.IsEmptyExceptNot() {
		return []byte("false"), nil
	}
	return util.ToJSONBytes(s.data, true), nil
}

func schemaFromJSON(msg []byte) (*Schema, error) {
	trimmed := bytes.TrimSpace(msg)
	if len(trimmed) == 0 {
		return nil, errors.New("empty JSON schema")
	}
	if len(trimmed) >= 4 {
		if bytes.Equal(trimmed, []byte("true")) {
			return trueSchemaData, nil
		}
		if bytes.Equal(trimmed, []byte("false")) {
			return falseSchemaData, nil
		}
	}
	if trimmed[0] == '[' {
		x, err := util.FromJSONObj[Schemas](msg)
		if err != nil {
			return nil, err
		}
		if len(x) == 0 {
			return nil, nil
		}
		if len(x) == 1 {
			return x[0], nil
		}
		return &Schema{data: data{dataValidations: dataValidations{dataArrayValidations: dataArrayValidations{PrefixItems: x}}}, bytes: msg}, nil
	}
	if trimmed[0] != '{' {
		return nil, errors.Errorf("invalid JSON schema root [%c], expected object or boolean", trimmed[0])
	}
	ret, err := util.FromJSONObj[data](msg)
	if err != nil {
		return nil, err
	}
	return &Schema{data: ret, bytes: msg}, nil
}
