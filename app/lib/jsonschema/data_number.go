package jsonschema

type DataNumber struct {
	MultipleOf       *float64 `json:"multipleOf,omitzero"`       // instance must be divisible by this number (strictly positive)
	Maximum          *float64 `json:"maximum,omitzero"`          // maximum inclusive value
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitzero"` // maximum exclusive value
	Minimum          *float64 `json:"minimum,omitzero"`          // minimum inclusive value
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitzero"` // minimum exclusive value
}

func (d DataNumber) IsEmpty() bool {
	return d.MultipleOf == nil && d.Maximum == nil && d.ExclusiveMaximum == nil && d.Minimum == nil && d.ExclusiveMinimum == nil
}

func (d DataNumber) Clone() DataNumber {
	return DataNumber{
		MultipleOf: d.MultipleOf, Maximum: d.Maximum, ExclusiveMaximum: d.ExclusiveMaximum,
		Minimum: d.Minimum, ExclusiveMinimum: d.ExclusiveMinimum,
	}
}
