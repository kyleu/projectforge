package jsonschema

func (s *Schema) HasDefinitions() bool {
	return !s.Defs.Empty() || !s.ExplicitDefs.Empty()
}

func (s *Schema) HasDescription() bool {
	return s.Description != ""
}

//nolint:gocyclo,cyclop
func (s *Schema) SchemaType() SchemaType {
	if s == nil {
		return SchemaTypeNull
	}
	if s.IsEmpty() {
		return SchemaTypeEmpty
	}

	switch s.Type {
	case "string":
		if len(s.Enum) > 0 {
			return SchemaTypeEnum
		}
		return SchemaTypeString
	case "integer":
		return SchemaTypeInteger
	case "object":
		if s.Ref != "" && s.Not == nil {
			return SchemaTypeRef
		}
		return SchemaTypeObject
	case "array":
		return SchemaTypeArray
	case nil:
		if s.Properties.Length() > 0 {
			return SchemaTypeObject
		}
		if len(s.Enum) > 0 {
			return SchemaTypeEnum
		}
	}
	if s.Ref != "" && s.Not == nil {
		return SchemaTypeRef
	}
	if s.Not != nil && s.Ref == "" {
		return SchemaTypeNot
	}
	return SchemaTypeObject
}
