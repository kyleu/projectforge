package jsonschema

import (
	"projectforge.dev/projectforge/app/util"
)

type dataValidations struct {
	// generic
	Type  any   `json:"type,omitzero"`  // expected data type(s) (string or array of "string", "number", "integer", "object", "array", "boolean", "null")
	Enum  []any `json:"enum,omitempty"` // array of allowed values
	Const any   `json:"const,omitzero"` // exact required value

	dataStringValidations
	dataNumberValidations
	dataArrayValidations
	dataObjectValidations
}

func (d dataValidations) IsEmpty() bool {
	return d.Type == nil && len(d.Enum) == 0 && d.Const == nil &&
		d.dataStringValidations.IsEmpty() && d.dataNumberValidations.IsEmpty() &&
		d.dataArrayValidations.IsEmpty() && d.dataObjectValidations.IsEmpty()
}

func (d dataValidations) Clone() dataValidations {
	return dataValidations{
		Type:                  d.Type,
		Enum:                  util.ArrayCopy(d.Enum),
		Const:                 d.Const,
		dataStringValidations: d.dataStringValidations.Clone(),
		dataNumberValidations: d.dataNumberValidations.Clone(),
		dataArrayValidations:  d.dataArrayValidations.Clone(),
		dataObjectValidations: d.dataObjectValidations.Clone(),
	}
}

type dataNumberValidations struct {
	MultipleOf       *float64 `json:"multipleOf,omitzero"`       // instance must be divisible by this number (strictly positive)
	Maximum          *float64 `json:"maximum,omitzero"`          // maximum inclusive value
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitzero"` // maximum exclusive value
	Minimum          *float64 `json:"minimum,omitzero"`          // minimum inclusive value
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitzero"` // minimum exclusive value
}

func (d dataNumberValidations) IsEmpty() bool {
	return d.MultipleOf == nil && d.Maximum == nil && d.ExclusiveMaximum == nil && d.Minimum == nil && d.ExclusiveMinimum == nil
}

func (d dataNumberValidations) Clone() dataNumberValidations {
	return dataNumberValidations{
		MultipleOf: d.MultipleOf, Maximum: d.Maximum, ExclusiveMaximum: d.ExclusiveMaximum,
		Minimum: d.Minimum, ExclusiveMinimum: d.ExclusiveMinimum,
	}
}

type dataStringValidations struct {
	MaxLength        *uint64 `json:"maxLength,omitzero"`        // maximum length (non-negative integer)
	MinLength        *uint64 `json:"minLength,omitzero"`        // minimum length (non-negative integer, default 0)
	Pattern          string  `json:"pattern,omitzero"`          // ecma 262 regular expression
	Format           string  `json:"format,omitzero"`           // predefined format (e.g., "date-time", "email", "ipv4")
	ContentEncoding  string  `json:"contentEncoding,omitzero"`  // encoding of the string content (e.g., "base64")
	ContentMediaType string  `json:"contentMediaType,omitzero"` // media type of the string content (e.g., "application/json")
	ContentSchema    *Schema `json:"contentSchema,omitzero"`    // schema for the decoded content if contentEncoding/contentMediaType are present
}

func (d dataStringValidations) IsEmpty() bool {
	return d.MaxLength == nil && d.MinLength == nil && d.Pattern == "" && d.Format == "" &&
		d.ContentEncoding == "" && d.ContentMediaType == "" && d.ContentSchema == nil
}

func (d dataStringValidations) Clone() dataStringValidations {
	return dataStringValidations{
		MaxLength: d.MaxLength, MinLength: d.MinLength, Pattern: d.Pattern, Format: d.Format,
		ContentEncoding: d.ContentEncoding, ContentMediaType: d.ContentMediaType, ContentSchema: d.ContentSchema.Clone(),
	}
}

type dataArrayValidations struct {
	Items            *Schema `json:"items,omitzero"`            // schema for array items (schema or boolean false). applied *after* prefixItems
	PrefixItems      Schemas `json:"prefixItems,omitzero"`      // array of schemas for tuple validation (items at specific indices)
	UnevaluatedItems *Schema `json:"unevaluatedItems,omitzero"` // validation for items not covered by `items` or `prefixItems` (schema or boolean)
	MaxItems         *uint64 `json:"maxItems,omitzero"`         // maximum number of items (non-negative integer)
	MinItems         *uint64 `json:"minItems,omitzero"`         // minimum number of items (non-negative integer, default 0)
	UniqueItems      *bool   `json:"uniqueItems,omitzero"`      // whether all items must be unique (default: false)
	Contains         *Schema `json:"contains,omitzero"`         // schema that at least one item must match
	MaxContains      *uint64 `json:"maxContains,omitzero"`      // maximum number of items matching `contains` (non-negative integer)
	MinContains      *uint64 `json:"minContains,omitzero"`      // minimum number of items matching `contains` (non-negative integer, default 1)
}

func (d dataArrayValidations) IsEmpty() bool {
	return d.Items == nil && len(d.PrefixItems) == 0 && d.UnevaluatedItems == nil && d.MaxItems == nil &&
		d.MinItems == nil && d.UniqueItems == nil && d.Contains == nil && d.MaxContains == nil && d.MinContains == nil
}

func (d dataArrayValidations) Clone() dataArrayValidations {
	return dataArrayValidations{
		Items: d.Items, PrefixItems: util.ArrayCopy(d.PrefixItems), UnevaluatedItems: d.UnevaluatedItems,
		MaxItems: d.MaxItems, MinItems: d.MinItems, UniqueItems: d.UniqueItems,
		Contains: d.Contains.Clone(), MaxContains: d.MaxContains, MinContains: d.MinContains,
	}
}

type dataObjectValidations struct {
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

func (d dataObjectValidations) IsEmpty() bool {
	return d.Properties.Empty() && d.PatternProperties.Empty() && d.AdditionalProperties == nil && !d.AllowTrailingCommas &&
		d.UnevaluatedProperties == nil && len(d.Required) == 0 && d.PropertyNames == nil && d.MaxProperties == nil &&
		d.MinProperties == nil && d.DependentRequired.Empty() && d.DependentSchemas.Empty()
}

func (d dataObjectValidations) Clone() dataObjectValidations {
	return dataObjectValidations{
		AdditionalProperties: d.AdditionalProperties.Clone(), AllowTrailingCommas: d.AllowTrailingCommas,
		UnevaluatedProperties: d.UnevaluatedProperties, Required: util.ArrayCopy(d.Required),
		PropertyNames: d.PropertyNames.Clone(), MaxProperties: d.MaxProperties, MinProperties: d.MinProperties,
		DependentRequired: d.DependentRequired.Clone(), DependentSchemas: d.DependentSchemas.Clone(),
	}
}
