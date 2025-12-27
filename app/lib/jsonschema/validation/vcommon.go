package validation

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
)

func validateCommon(s *jsonschema.Schema) error {
	if len(s.Required) > 0 && (s.Properties.Empty()) {
		return errors.Errorf("schema [%s] has required fields, but no properties", s.String())
	}
	return nil
}
