package jsonschema

import "projectforge.dev/projectforge/app/util"

type DataValidations struct {
	// generic
	Type  any   `json:"type,omitzero"`  // expected data type(s) (string or array of "string", "number", "integer", "object", "array", "boolean", "null")
	Enum  []any `json:"enum,omitempty"` // array of allowed values
	Const any   `json:"const,omitzero"` // exact required value
}

func (d DataValidations) IsEmpty() bool {
	return d.Type == nil && len(d.Enum) == 0 && d.Const == nil
}

func (d DataValidations) Clone() DataValidations {
	return DataValidations{
		Type:  d.Type,
		Enum:  util.ArrayCopy(d.Enum),
		Const: d.Const,
	}
}
