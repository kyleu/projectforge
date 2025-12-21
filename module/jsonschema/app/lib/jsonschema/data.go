package jsonschema

import (
	"encoding/json/jsontext"

	"{{{ .Package }}}/app/util"
)

type data struct {
	dataCore
	dataAnnotations
	dataValidations
	dataApplicators
	Examples []any                     `json:"examples,omitempty"` // array of example values that validate against the schema
	Unknown  map[string]jsontext.Value `json:",unknown"`           // extra fields, usually metadata
}

func (d data) IsEmpty() bool {
	return d.dataCore.IsEmpty() && d.dataAnnotations.IsEmpty() && d.dataValidations.IsEmpty() && d.dataApplicators.IsEmpty() &&
		len(d.Examples) == 0 && len(d.Unknown) == 0
}

func (d data) IsEmptyExceptNot() bool {
	return d.Not != nil && d.dataCore.IsEmpty() && d.dataAnnotations.IsEmpty() &&
		d.dataValidations.IsEmpty() && d.dataApplicators.IsEmptyExceptNot() &&
		len(d.Examples) == 0 && len(d.Unknown) == 0
}

func (d data) Clone() data {
	return data{
		dataCore:        d.dataCore.Clone(),
		dataAnnotations: d.dataAnnotations.Clone(),
		dataValidations: d.dataValidations.Clone(),
		dataApplicators: d.dataApplicators.Clone(),
		Examples:        util.ArrayCopy(d.Examples),
		Unknown:         util.MapClone(d.Unknown),
	}
}
