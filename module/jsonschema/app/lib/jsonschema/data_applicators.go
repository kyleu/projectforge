package jsonschema

type dataApplicators struct {
	// conditional
	If   *Schema `json:"if,omitzero"`   // if this schema validates, `then` must also validate
	Then *Schema `json:"then,omitzero"` // schema applied if `if` validates
	Else *Schema `json:"else,omitzero"` // schema applied if `if` does not validate

	// boolean logic
	AllOf Schemas `json:"allOf,omitempty"` // instance must validate against all of these schemas
	AnyOf Schemas `json:"anyOf,omitempty"` // instance must validate against at least one of these schemas
	OneOf Schemas `json:"oneOf,omitempty"` // instance must validate against exactly one of these schemas
	Not   *Schema `json:"not,omitempty"`   // instance must not validate against this schema
}

func (d dataApplicators) IsEmpty() bool {
	return d.If == nil && d.Then == nil && d.Else == nil && len(d.AllOf) == 0 && len(d.AnyOf) == 0 && len(d.OneOf) == 0 && d.Not == nil
}

func (d dataApplicators) IsEmptyExceptNot() bool {
	return d.If == nil && d.Then == nil && d.Else == nil && len(d.AllOf) == 0 && len(d.AnyOf) == 0 && len(d.OneOf) == 0
}

func (d dataApplicators) Clone() dataApplicators {
	return dataApplicators{
		If: d.If.Clone(), Then: d.Then.Clone(), Else: d.Else.Clone(),
		AllOf: d.AllOf.Clone(), AnyOf: d.AnyOf.Clone(), OneOf: d.OneOf.Clone(), Not: d.Not.Clone(),
	}
}
