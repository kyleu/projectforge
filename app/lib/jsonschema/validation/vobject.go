package validation

import (
	"slices"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
)

func validateObject(s *jsonschema.Schema) error {
	for _, x := range s.Required {
		if !slices.Contains(s.Properties.Order, x) {
			return errors.Errorf("schema [%s] has required field [%s], but isn't in properties", s.String(), x)
		}
	}
	return nil
}

func validateNonObject(s *jsonschema.Schema) error {
	if len(s.Required) > 0 {
		return errors.Errorf("schema [%s] has required fields, but isn't an object", s.String())
	}
	return nil
}
