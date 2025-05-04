package jsonschema

import (
	"slices"

	"github.com/pkg/errors"
)

func (s *Schema) Validate() error {
	if err := validateType(s); err != nil {
		return err
	}
	if err := validateRequired(s); err != nil {
		return err
	}
	return nil
}

func validateType(s *Schema) error {
	switch t := s.Type.(type) {
	case string:
		// noop
	default:
		return errors.Errorf("schema [%s] type [%T] can't be handled", s.String(), t)
	}
	return nil
}

func validateRequired(s *Schema) error {
	if s.Type != "object" && len(s.Required) > 0 {
		return errors.Errorf("schema [%s] has required fields, but isn't an object", s.String())
	}
	for _, x := range s.Required {
		if !slices.Contains(s.Properties.Order, x) {
			return errors.Errorf("schema [%s] has required field [%s], but isn't in properties", s.String(), x)
		}
	}
	return nil
}
