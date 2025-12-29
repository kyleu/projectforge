package validation

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
)

func validateEnum(s *jsonschema.Schema) error {
	if len(s.Enum) > 0 && (!s.Properties.Empty() || len(s.Required) > 0) {
		return errors.Errorf("schema [%s] is an enum, but has properties or required fields", s.String())
	}
	return nil
}

func validateNonEnum(s *jsonschema.Schema) error {
	if len(s.Enum) > 0 {
		return errors.Errorf("schema [%s] is not an enum, but has enum values", s.String())
	}
	return nil
}
