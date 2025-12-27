package jsonschema

import "{{{ .Package }}}/app/util"

type DataObject struct {
	Properties            *util.OrderedMap[*Schema]  `json:"properties,omitzero"`            // schemas for named properties
	PatternProperties     *util.OrderedMap[*Schema]  `json:"patternProperties,omitzero"`     // schemas for properties matching regex patterns
	AdditionalProperties  *Schema                    `json:"additionalProperties,omitzero"`  // schema or boolean
	AllowTrailingCommas   bool                       `json:"allowTrailingCommas,omitzero"`   // indicates that trailing commas are allowed
	UnevaluatedProperties any                        `json:"unevaluatedProperties,omitzero"` // validation for properties not covered by `properties`
	Required              []string                   `json:"required,omitempty"`             // array of required property names
	PropertyNames         *Schema                    `json:"propertyNames,omitzero"`         // schema for property names
	MaxProperties         *uint64                    `json:"maxProperties,omitzero"`         // maximum number of properties (non-negative integer)
	MinProperties         *uint64                    `json:"minProperties,omitzero"`         // minimum number of properties (non-negative integer, default 0)
	DependentRequired     *util.OrderedMap[[]string] `json:"dependentRequired,omitzero"`     // properties required based on the presence of other properties
	DependentSchemas      *util.OrderedMap[*Schema]  `json:"dependentSchemas,omitzero"`      // schemas applied based on the presence of other properties
}

func (d DataObject) IsEmpty() bool {
	return d.Properties.Empty() && d.PatternProperties.Empty() && d.AdditionalProperties == nil && !d.AllowTrailingCommas &&
		d.UnevaluatedProperties == nil && len(d.Required) == 0 && d.PropertyNames == nil && d.MaxProperties == nil &&
		d.MinProperties == nil && d.DependentRequired.Empty() && d.DependentSchemas.Empty()
}

func (d DataObject) Clone() DataObject {
	return DataObject{
		Properties: d.Properties.Clone(), PatternProperties: d.PatternProperties.Clone(),
		AdditionalProperties: d.AdditionalProperties.Clone(), AllowTrailingCommas: d.AllowTrailingCommas,
		UnevaluatedProperties: d.UnevaluatedProperties, Required: util.ArrayCopy(d.Required),
		PropertyNames: d.PropertyNames.Clone(), MaxProperties: d.MaxProperties, MinProperties: d.MinProperties,
		DependentRequired: d.DependentRequired.Clone(), DependentSchemas: d.DependentSchemas.Clone(),
	}
}
