package jsonschema

import (
	"encoding/json/jsontext"

	"projectforge.dev/projectforge/app/util"
)

type Data struct {
	DataCore
	DataAnnotations
	DataValidations
	DataString
	DataNumber
	DataArray
	DataObject
	DataApplicators
	Examples []any                     `json:"examples,omitempty"` // array of example values that validate against the schema
	Unknown  map[string]jsontext.Value `json:",unknown"`           // extra fields, usually metadata
}

func (d Data) IsEmpty() bool {
	return d.DataCore.IsEmpty() && d.DataAnnotations.IsEmpty() && d.DataValidations.IsEmpty() &&
		d.DataString.IsEmpty() && d.DataNumber.IsEmpty() && d.DataArray.IsEmpty() && d.DataObject.IsEmpty() &&
		d.DataApplicators.IsEmpty() && len(d.Examples) == 0 && len(d.Unknown) == 0
}

func (d Data) IsEmptyExceptNot() bool {
	return d.Not != nil && d.DataCore.IsEmpty() && d.DataAnnotations.IsEmpty() && d.DataValidations.IsEmpty() &&
		d.DataString.IsEmpty() && d.DataNumber.IsEmpty() && d.DataArray.IsEmpty() && d.DataObject.IsEmpty() &&
		d.DataApplicators.IsEmptyExceptNot() && len(d.Examples) == 0 && len(d.Unknown) == 0
}

func (d Data) Clone() Data {
	return Data{
		DataCore:        d.DataCore.Clone(),
		DataAnnotations: d.DataAnnotations.Clone(),
		DataValidations: d.DataValidations.Clone(),
		DataString:      d.DataString.Clone(),
		DataNumber:      d.DataNumber.Clone(),
		DataArray:       d.DataArray.Clone(),
		DataObject:      d.DataObject.Clone(),
		DataApplicators: d.DataApplicators.Clone(),
		Examples:        util.ArrayCopy(d.Examples),
		Unknown:         util.MapClone(d.Unknown),
	}
}
