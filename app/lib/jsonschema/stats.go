package jsonschema

func (s *Schema) HasDefinitions() bool {
	return !s.Defs.Empty() || !s.ExplicitDefs.Empty()
}

func (s *Schema) HasDescription() bool {
	return s.Description != ""
}
