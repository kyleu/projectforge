package validation

import (
	"projectforge.dev/projectforge/app/lib/jsonschema"
)

func ValidateCollection(c *jsonschema.Collection) error {
	for _, sch := range c.SchemaMap {
		if _, err := Validate(sch); err != nil {
			return err
		}
	}
	return nil
}

func Validate(s *jsonschema.Schema) (jsonschema.SchemaType, error) {
	t := s.DetectSchemaType()
	if err := validateCommon(s); err != nil {
		return t, err
	}
	fn := func(st jsonschema.SchemaType, y func(*jsonschema.Schema) error, n func(*jsonschema.Schema) error) error {
		if t == st {
			return y(s)
		}
		return n(s)
	}
	if err := fn(t, validateArray, validateNonArray); err != nil {
		return t, err
	}
	if err := fn(t, validateBoolean, validateNonBoolean); err != nil {
		return t, err
	}
	if err := fn(t, validateEnum, validateNonEnum); err != nil {
		return t, err
	}
	if err := fn(t, validateInt, validateNonInt); err != nil {
		return t, err
	}
	if err := fn(t, validateNot, validateNonNot); err != nil {
		return t, err
	}
	if err := fn(t, validateNull, validateNonNull); err != nil {
		return t, err
	}
	if err := fn(t, validateNumber, validateNonNumber); err != nil {
		return t, err
	}
	if err := fn(t, validateObject, validateNonObject); err != nil {
		return t, err
	}
	if err := fn(t, validateRef, validateNonRef); err != nil {
		return t, err
	}
	if err := fn(t, validateString, validateNonString); err != nil {
		return t, err
	}
	if err := fn(t, validateUnion, validateNonUnion); err != nil {
		return t, err
	}
	return t, nil
}
