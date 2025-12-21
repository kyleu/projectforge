package jsonschema

import (
	"slices"

	"github.com/pkg/errors"
)

func (s *Schema) Validate() error {
	t := s.DetectSchemaType()
	if t == SchemaTypeObject {

	} else {

	}
	if err := validateRequired(s, t); err != nil {
		return err
	}
	return nil
}

func validateRequired(s *Schema, t SchemaType) error {
	if t != SchemaTypeObject && len(s.Required) > 0 {
		return errors.Errorf("schema [%s] has required fields, but isn't an object", s.String())
	}
	if len(s.Required) > 0 && (s.Properties.Empty()) {
		return errors.Errorf("schema [%s] has required fields, but no properties", s.String())
	}
	for _, x := range s.Required {
		if !slices.Contains(s.Properties.Order, x) {
			return errors.Errorf("schema [%s] has required field [%s], but isn't in properties", s.String(), x)
		}
	}
	return nil
}

func (s *Schema) DetectSchemaType() SchemaType {
	if s == nil {
		return SchemaTypeNull
	}
	switch s.Type {
	case "string":
		if len(s.Enum) > 0 {
			return SchemaTypeEnum
		}
		return SchemaTypeString
	case "integer":
		return SchemaTypeInteger
	case "array":
		return SchemaTypeArray
	case "number":
		return SchemaTypeNumber
	case "boolean":
		return SchemaTypeBoolean
	case "null":
		return SchemaTypeNull
	default:
		if !s.Properties.Empty() {
			return SchemaTypeObject
		}
		if len(s.Required) > 0 {
			return SchemaTypeObject
		}
		if len(s.Enum) > 0 {
			return SchemaTypeEnum
		}
		if s.Ref != "" && s.Not == nil {
			return SchemaTypeRef
		}
		if s.Not != nil && s.Ref == "" {
			return SchemaTypeNot
		}
		if len(s.OneOf) > 0 || len(s.AnyOf) > 0 && len(s.AllOf) > 0 {
			// return SchemaTypeUnion
		}
		return SchemaTypeUnknown
	}
}
