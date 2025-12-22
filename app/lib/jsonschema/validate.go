package jsonschema

import (
	"slices"

	"github.com/pkg/errors"
)

func (s *Schema) DetectSchemaType() SchemaType {
	ret, _ := s.Validate()
	return ret
}

func (s *Schema) Validate() (SchemaType, error) {
	if s == nil {
		return SchemaTypeNull, nil
	}
	sendErr := func(msg string, args ...any) (SchemaType, error) {
		return SchemaTypeUnknown, errors.Errorf(msg, args...)
	}
	if len(s.Required) > 0 && (s.Properties.Empty()) {
		return sendErr("schema [%s] has required fields, but no properties", s.String())
	}

	switch s.Type {
	case "string":
		if len(s.Enum) > 0 {
			return SchemaTypeEnum, validateEnum(s)
		}
		return SchemaTypeString, validateString(s)
	case "integer":
		return SchemaTypeInteger, validateInt(s)
	case "array":
		return SchemaTypeArray, validateArray(s)
	case "number":
		return SchemaTypeNumber, validateNumber(s)
	case "boolean":
		return SchemaTypeBoolean, validateBoolean(s)
	case "null":
		return SchemaTypeNull, validateNull(s)
	default:
		if s.HasProperties() {
			return SchemaTypeObject, validateObject(s)
		}
		if len(s.Enum) > 0 {
			return SchemaTypeEnum, validateEnum(s)
		}
		if s.Ref != "" && s.Not == nil {
			return SchemaTypeRef, validateRef(s)
		}
		if s.Not != nil && s.Ref == "" {
			return SchemaTypeNot, validateNot(s)
		}
		if len(s.OneOf) > 0 || len(s.AnyOf) > 0 && len(s.AllOf) > 0 {
			return SchemaTypeUnion, validateUnion(s)
		}
		if s.IsEmpty() {
			return SchemaTypeEmpty, nil
		}
		if s.IsEmptyExceptNot() {
			return SchemaTypeNot, nil
		}
		return SchemaTypeUnknown, nil
	}
}

func validateObject(s *Schema) error {
	for _, x := range s.Required {
		if !slices.Contains(s.Properties.Order, x) {
			return errors.Errorf("schema [%s] has required field [%s], but isn't in properties", s.String(), x)
		}
	}
	return nil
}

func validateNonObject(s *Schema) error {
	if len(s.Required) > 0 {
		return errors.Errorf("schema [%s] has required fields, but isn't an object", s.String())
	}
	return nil
}

func validateEnum(s *Schema) error {
	if len(s.Enum) > 0 && (!s.Properties.Empty() || len(s.Required) > 0) {
		return errors.Errorf("schema [%s] is an enum, but has properties or required fields", s.String())
	}
	return validateNonObject(s)
}

func validateString(s *Schema) error {
	return validateNonObject(s)
}

func validateInt(s *Schema) error {
	return validateNonObject(s)
}

func validateArray(s *Schema) error {
	return validateNonObject(s)
}

func validateNumber(s *Schema) error {
	return validateNonObject(s)
}

func validateBoolean(s *Schema) error {
	return validateNonObject(s)
}

func validateNull(s *Schema) error {
	return validateNonObject(s)
}

func validateRef(s *Schema) error {
	return validateNonObject(s)
}

func validateNot(s *Schema) error {
	return validateNonObject(s)
}

func validateUnion(s *Schema) error {
	return validateNonObject(s)
}
